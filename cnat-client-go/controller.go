package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	cnatv1alpha1 "github.com/ZhengHe-MD/programming-kubernetes/cnat-client-go/pkg/apis/cnat/v1alpha1"
	clientset "github.com/ZhengHe-MD/programming-kubernetes/cnat-client-go/pkg/generated/clientset/versioned"
	"github.com/ZhengHe-MD/programming-kubernetes/cnat-client-go/pkg/generated/informers/externalversions/cnat/v1alpha1"
	listers "github.com/ZhengHe-MD/programming-kubernetes/cnat-client-go/pkg/generated/listers/cnat/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1informer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	corev1lister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

type Controller struct {
	kubeClientset kubernetes.Interface
	cnatClientset clientset.Interface

	atLister listers.AtLister
	atSynced cache.InformerSynced

	podLister corev1lister.PodLister
	podSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
}

func NewController(
	kubeClientset kubernetes.Interface,
	cnatClientset clientset.Interface,
	atInformer v1alpha1.AtInformer,
	podInformer corev1informer.PodInformer) *Controller {

	utilruntime.Must(cnatv1alpha1.AddToScheme(scheme.Scheme))

	controller := &Controller{
		kubeClientset: kubeClientset,
		cnatClientset: cnatClientset,
		atLister:      atInformer.Lister(),
		atSynced:      atInformer.Informer().HasSynced,
		podLister:     podInformer.Lister(),
		podSynced:     podInformer.Informer().HasSynced,
		workqueue: workqueue.NewRateLimitingQueueWithConfig(
			workqueue.DefaultControllerRateLimiter(),
			workqueue.RateLimitingQueueConfig{Name: "Ats"}),
	}

	klog.Info("Setting up event handlers")
	_, _ = atInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueAt,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueAt(new)
		},
		DeleteFunc: nil,
	})

	_, _ = podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueuePod,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueuePod(new)
		},
		DeleteFunc: nil,
	})

	return controller
}

func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("Starting cnat client-go controller")

	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.atSynced); !ok {
		return fmt.Errorf("failed to wait for At caches to sync")
	}

	if ok := cache.WaitForCacheSync(stopCh, c.podSynced); !ok {
		return fmt.Errorf("failed to wait for pod caches to sync")
	}

	klog.Info("Starting workers")
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}

	// wrap this block in a func so we can defer c.workerqueue.Done
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)

		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		when, err := c.syncHandler(key)
		if err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}

		if when != time.Duration(0) {
			c.workqueue.AddAfter(key, when)
		} else {
			c.workqueue.Forget(obj)
		}
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) (time.Duration, error) {
	klog.Infof("=== Reconciling At %s", key)
	ctx := context.TODO()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return time.Duration(0), nil
	}

	original, err := c.atLister.Ats(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("at '%s' in work queue no longer exists", key))
			return time.Duration(0), nil
		}
		return time.Duration(0), err
	}

	// Clone because the original object is owned by the lister
	instance := original.DeepCopy()

	if instance.Status.Phase == "" {
		instance.Status.Phase = cnatv1alpha1.PhasePending
	}

	// PENDING -> RUNNING -> DONE
	switch instance.Status.Phase {
	case cnatv1alpha1.PhasePending:
		klog.Infof("instance %s: phase=PENDING", key)
		klog.Infof("instance %s: checking schedule %q", key, instance.Spec.Schedule)
		d, err := timeUntilSchedule(instance.Spec.Schedule)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("schedule parsing failed: %v", err))
			// Error reading the schedule, wait until it is fixed.
			return time.Duration(0), err
		}
		klog.Infof("instance %s: schedule parsing done: diff=%v", key, d)
		if d > 0 {
			return d, nil
		}

		klog.Infof("instance %s: it's time! Ready to execute: %s", key, instance.Spec.Command)
		instance.Status.Phase = cnatv1alpha1.PhaseRunning
	case cnatv1alpha1.PhaseRunning:
		klog.Infof("instance %s: Phase=RUNNING", key)
		pod := newPodForCR(instance)

		// Set At instance as the owner and controller
		owner := metav1.NewControllerRef(instance, cnatv1alpha1.SchemeGroupVersion.WithKind("At"))
		pod.ObjectMeta.OwnerReferences = append(pod.ObjectMeta.OwnerReferences, *owner)

		// Try to see if the pod already exists, if not create a new one
		found, err := c.kubeClientset.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			found, err = c.kubeClientset.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
			if err != nil {
				return time.Duration(0), err
			}
			klog.Infof("instance %s: pod launched: name=%s", key, pod.Name)
		} else if err != nil {
			return time.Duration(0), err
		} else if found.Status.Phase == corev1.PodFailed || found.Status.Phase == corev1.PodSucceeded {
			klog.Infof("instance %s: container terminated: reason=%q message=%q", key, found.Status.Reason, found.Status.Message)
			instance.Status.Phase = cnatv1alpha1.PhaseDone
		} else {
			// pod is pending, running or in an unknown state, don't requeue
			// in this case because we'll be called again when the pod status changes.
			return time.Duration(0), nil
		}
	case cnatv1alpha1.PhaseDone:
		klog.Infof("instance %s: phase=DONE", key)
		return time.Duration(0), nil
	default:
		klog.Infof("instance %s: NOP", key)
		return time.Duration(0), nil
	}

	if !reflect.DeepEqual(original, instance) {
		_, err = c.cnatClientset.CnatV1alpha1().Ats(instance.Namespace).UpdateStatus(context.TODO(), instance, metav1.UpdateOptions{})
		if err != nil {
			return time.Duration(0), err
		}
	}

	return time.Duration(0), nil
}

func (c *Controller) enqueueAt(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

func (c *Controller) enqueuePod(obj interface{}) {
	var pod *corev1.Pod
	var ok bool
	if pod, ok = obj.(*corev1.Pod); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding pod, invalid type"))
			return
		}
		pod, ok = tombstone.Obj.(*corev1.Pod)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding pod tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted pod '%s' from tombstone", pod.GetName())
	}
	if ownerRef := metav1.GetControllerOf(pod); ownerRef != nil {
		if ownerRef.Kind != "At" {
			return
		}

		at, err := c.atLister.Ats(pod.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned pod '%s' of At '%s'",
				pod.GetName(), ownerRef.Name)
			return
		}

		klog.Infof("enqueuing At %s/%s because pod changed", at.Namespace, at.Name)
		c.enqueueAt(at)
	}
}

func newPodForCR(cr *cnatv1alpha1.At) *corev1.Pod {
	labels := map[string]string{"app": cr.Name}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: strings.Split(cr.Spec.Command, " "),
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}

func timeUntilSchedule(schedule string) (time.Duration, error) {
	s, err := time.Parse("2006-01-02T15:04:05Z", schedule)
	if err != nil {
		return time.Duration(0), err
	}
	return s.Sub(time.Now()), nil
}
