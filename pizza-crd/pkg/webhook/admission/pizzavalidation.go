package admission

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ZhengHe-MD/programming-kubernetes/pizza-crd/pkg/apis/restaurant/install"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-crd/pkg/apis/restaurant/v1alpha1"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-crd/pkg/apis/restaurant/v1beta1"
	restaurantinformers "github.com/ZhengHe-MD/programming-kubernetes/pizza-crd/pkg/generated/informers/externalversions"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func init() {
	utilruntime.Must(admissionv1beta1.AddToScheme(scheme))
	install.Install(scheme)
}

func ServePizzaValidation(informers restaurantinformers.SharedInformerFactory) func(http.ResponseWriter, *http.Request) {
	toppingInformer := informers.Restaurant().V1alpha1().Toppings().Informer()
	toppingLister := informers.Restaurant().V1alpha1().Toppings().Lister()

	return func(w http.ResponseWriter, req *http.Request) {
		if !toppingInformer.HasSynced() {
			responsewriters.InternalError(w, req, fmt.Errorf("informers not ready"))
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			responsewriters.InternalError(w, req, fmt.Errorf("failed to read body: %v", err))
			return
		}

		gv := admissionv1beta1.SchemeGroupVersion
		reviewGVK := gv.WithKind("AdmissionReview")
		obj, gvk, err := codecs.UniversalDeserializer().Decode(body, &reviewGVK, &admissionv1beta1.AdmissionReview{})
		if err != nil {
			responsewriters.InternalError(w, req, fmt.Errorf("failed to decode body: %v", err))
			return
		}
		review, ok := obj.(*admissionv1beta1.AdmissionReview)
		if !ok {
			responsewriters.InternalError(w, req, fmt.Errorf("unexpected GroupVersionKind: %s", gvk))
			return
		}
		if review.Request == nil {
			responsewriters.InternalError(w, req, fmt.Errorf("malformed admission review: request is nil"))
			return
		}
		review.Response = &admissionv1beta1.AdmissionResponse{UID: review.Request.UID}

		if review.Request.Object.Object == nil {
			var err error
			review.Request.Object.Object, _, err = codecs.UniversalDeserializer().Decode(review.Request.Object.Raw, nil, nil)
			if err != nil {
				review.Response.Result = &metav1.Status{
					Message: err.Error(),
					Status:  metav1.StatusFailure,
				}
				responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, gvk.GroupVersion(), w, req, http.StatusOK, review, false)
				return
			}
		}

		switch pizza := review.Request.Object.Object.(type) {
		case *v1alpha1.Pizza:
			for _, topping := range pizza.Spec.Toppings {
				if _, err := toppingLister.Get(topping); err != nil && !errors.IsNotFound(err) {
					responsewriters.InternalError(w, req, fmt.Errorf("failed to lookup topping %q: %v", topping, err))
					return
				} else if errors.IsNotFound(err) {
					review.Response.Result = &metav1.Status{
						Message: fmt.Sprintf("topping %q not known", topping),
						Status:  metav1.StatusFailure,
					}
					responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, gvk.GroupVersion(), w, req, http.StatusOK, review, false)
					return
				}
			}
			review.Response.Allowed = true
		case *v1beta1.Pizza:
			for _, topping := range pizza.Spec.Toppings {
				if _, err := toppingLister.Get(topping.Name); err != nil && !errors.IsNotFound(err) {
					responsewriters.InternalError(w, req, fmt.Errorf("failed to lookup topping %q: %v", topping, err))
					return
				} else if errors.IsNotFound(err) {
					review.Response.Result = &metav1.Status{
						Message: fmt.Sprintf("topping %q not known", topping),
						Status:  metav1.StatusFailure,
					}
					responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, gvk.GroupVersion(), w, req, http.StatusOK, review, false)
					return
				}
			}
			review.Response.Allowed = true
		default:
			review.Response.Result = &metav1.Status{
				Message: fmt.Sprintf("unexpected type %T", review.Request.Object.Object),
				Status:  metav1.StatusFailure,
			}
		}
		responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, gvk.GroupVersion(), w, req, http.StatusOK, review, false)
	}
}
