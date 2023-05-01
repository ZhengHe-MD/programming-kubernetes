/*
Copyright 2023 ZhengHe-MD.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cnatv1alpha1 "github.com/ZhengHe-MD/cnat-operator-sdk/api/v1alpha1"
)

// AtReconciler reconciles a At object
type AtReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cnat.programming-kubernetes.info,resources=ats,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cnat.programming-kubernetes.info,resources=ats/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cnat.programming-kubernetes.info,resources=ats/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the At object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *AtReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx).WithValues("namespace", req.Namespace, "at", req.Name)
	reqLogger.Info("=== Reconciling At")
	// Fetch the At instance
	instance := &cnatv1alpha1.At{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "error reading the At instance, requeue the request")
		return ctrl.Result{}, err
	}

	// If no phase set, default to pending (the initial phase):
	if instance.Status.Phase == "" {
		instance.Status.Phase = cnatv1alpha1.PhasePending
	}

	// PENDING -> RUNNING -> DONE
	switch instance.Status.Phase {
	case cnatv1alpha1.PhasePending:
		reqLogger.Info("phase: PENDING")
		reqLogger.Info("checking schedule", "target", instance.Spec.Schedule)
		d, err := timeUntilSchedule(instance.Spec.Schedule)
		if err != nil {
			reqLogger.Error(err, "schedule parsing failure")
			// Error reading the schedule. Wait until it is fixed.
			return ctrl.Result{}, err
		}
		reqLogger.Info("schedule parsing done", "result", fmt.Sprintf("diff=%v", d))
		if d > 0 {
			return ctrl.Result{RequeueAfter: d}, nil
		}
		reqLogger.Info("ready to execute", "command", instance.Spec.Command)
		instance.Status.Phase = cnatv1alpha1.PhaseRunning
	case cnatv1alpha1.PhaseRunning:
		reqLogger.Info("phase: RUNNING")
		pod := newPodForCR(instance)
		// Set At instance as the owner and controller
		if err = ctrl.SetControllerReference(instance, pod, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		found := &corev1.Pod{}
		err = r.Get(ctx, client.ObjectKeyFromObject(pod), found)
		if err != nil && errors.IsNotFound(err) {
			err = r.Create(ctx, pod)
			if err != nil {
				return ctrl.Result{}, err
			}
			reqLogger.Info("pod launched", "name", pod.Name)
		} else if err != nil {
			return ctrl.Result{}, err
		} else if found.Status.Phase == corev1.PodFailed || found.Status.Phase == corev1.PodSucceeded {
			reqLogger.Info("container terminated",
				"reason", found.Status.Reason,
				"message", found.Status.Message)
			instance.Status.Phase = cnatv1alpha1.PhaseDone
		} else {
			// pod is pending, running or in an unknown state, don't requeue
			// in this case because we'll be called again when the pod status changes.
			return ctrl.Result{}, nil
		}
	case cnatv1alpha1.PhaseDone:
		reqLogger.Info("phase: DONE")
		return ctrl.Result{}, nil
	default:
		reqLogger.Info("NOP")
		return ctrl.Result{}, nil
	}

	// Update the At instance, setting the status to the respective phase
	if err = r.Status().Update(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}

	// Don't requeue - we will reconcile again when there's a change to the object.
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AtReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cnatv1alpha1.At{}).
		Watches(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
			OwnerType:    &cnatv1alpha1.At{},
			IsController: true,
		}).
		Complete(r)
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
