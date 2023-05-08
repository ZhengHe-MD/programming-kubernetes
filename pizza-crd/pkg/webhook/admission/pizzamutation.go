package admission

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ZhengHe-MD/programming-kubernetes/pizza-crd/pkg/apis/restaurant/v1alpha1"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-crd/pkg/apis/restaurant/v1beta1"
	"github.com/appscode/jsonpatch"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/klog/v2"
)

func ServePizzaAdmit(w http.ResponseWriter, req *http.Request) {
	// read body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		responsewriters.InternalError(w, req, fmt.Errorf("failed to read body: %v", err))
		return
	}

	reviewGVK := admissionv1beta1.SchemeGroupVersion.WithKind("AdmissionReview")
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
		responsewriters.InternalError(w, req, fmt.Errorf("unexpected nil request"))
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

	var bs []byte
	switch pizza := review.Request.Object.Object.(type) {
	case *v1alpha1.Pizza:
		if len(pizza.Spec.Toppings) == 0 {
			pizza.Spec.Toppings = []string{"tomato", "mozzarella", "salami"}
		}
		bs, err = json.Marshal(pizza)
		if err != nil {
			responsewriters.InternalError(w, req, fmt.Errorf("unexpected encoding error: %v", err))
			return
		}
	case *v1beta1.Pizza:
		if len(pizza.Spec.Toppings) == 0 {
			pizza.Spec.Toppings = []v1beta1.PizzaTopping{
				{"tomato", 1},
				{"mozzarella", 1},
				{"salami", 1},
			}
		}
		bs, err = json.Marshal(pizza)
		if err != nil {
			responsewriters.InternalError(w, req, fmt.Errorf("unexpected encoding error: %v", err))
			return
		}
	default:
		review.Response.Result = &metav1.Status{
			Message: fmt.Sprintf("unexpected type %T", review.Request.Object.Object),
			Status:  metav1.StatusFailure,
		}
		responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, gvk.GroupVersion(), w, req, http.StatusOK, review, false)
		return
	}

	klog.V(2).Infof("Defaulting %s/%s in version %s", review.Request.Namespace, review.Request.Name, gvk)

	// compare original and defaulted version
	orig := review.Request.Object.Raw

	ops, err := jsonpatch.CreatePatch(orig, bs)
	if err != nil {
		responsewriters.InternalError(w, req, fmt.Errorf("unexpected diff error: %v", err))
		return
	}
	review.Response.Patch, err = json.Marshal(ops)
	if err != nil {
		responsewriters.InternalError(w, req, fmt.Errorf("unexpected patch encoding error: %v", err))
		return
	}

	typ := admissionv1beta1.PatchTypeJSONPatch
	review.Response.PatchType = &typ
	review.Response.Allowed = true

	responsewriters.WriteObjectNegotiated(codecs, negotiation.DefaultEndpointRestrictions, gvk.GroupVersion(), w, req, http.StatusOK, review, false)
}
