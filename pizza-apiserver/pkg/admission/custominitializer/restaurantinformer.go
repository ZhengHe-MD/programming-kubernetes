package custominitializer

import (
	informers "github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/generated/informers/externalversions"
	"k8s.io/apiserver/pkg/admission"
)

var _ admission.PluginInitializer = restaurantInformerPluginInitializer{}

type restaurantInformerPluginInitializer struct {
	informers informers.SharedInformerFactory
}

func New(sharedInformerFactory informers.SharedInformerFactory) restaurantInformerPluginInitializer {
	return restaurantInformerPluginInitializer{informers: sharedInformerFactory}
}

func (r restaurantInformerPluginInitializer) Initialize(plugin admission.Interface) {
	if wants, ok := plugin.(WantsRestaurantInformerFactory); ok {
		wants.SetRestaurantInformerFactory(r.informers)
	}
}
