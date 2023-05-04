package custominitializer

import (
	informers "github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/generated/informers/externalversions"
	"k8s.io/apiserver/pkg/admission"
)

type WantsRestaurantInformerFactory interface {
	SetRestaurantInformerFactory(factory informers.SharedInformerFactory)
	admission.InitializationValidator
}
