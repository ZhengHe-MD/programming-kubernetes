package pizzatoppings

import (
	"context"
	"fmt"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/admission/custominitializer"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant"
	informers "github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/generated/informers/externalversions"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/generated/listers/restaurant/v1alpha1"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"
)

func Register(plugins *admission.Plugins) {
	plugins.Register("PizzaToppings", func(config io.Reader) (admission.Interface, error) {
		return New()
	})
}

var _ custominitializer.WantsRestaurantInformerFactory = (*PizzaToppingPlugin)(nil)
var _ admission.ValidationInterface = (*PizzaToppingPlugin)(nil)

type PizzaToppingPlugin struct {
	*admission.Handler
	toppingLister v1alpha1.ToppingLister
}

func New() (*PizzaToppingPlugin, error) {
	return &PizzaToppingPlugin{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}, nil
}

func (p *PizzaToppingPlugin) Validate(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	// we are only interested in pizzas
	if a.GetKind().GroupKind() != restaurant.Kind("Pizza") {
		return nil
	}

	if !p.WaitForReady() {
		return admission.NewForbidden(a, fmt.Errorf("not yet ready to handle request"))
	}

	pizza := a.GetObject().(*restaurant.Pizza)
	for _, top := range pizza.Spec.Toppings {
		if _, err := p.toppingLister.Get(top.Name); err != nil && errors.IsNotFound(err) {
			return admission.NewForbidden(a, fmt.Errorf("unknown topping: %s", top.Name))
		} else if err != nil {
			return admission.NewForbidden(a, fmt.Errorf("internal server error while checking topping: %s", top.Name))
		}
	}
	return nil
}

func (p *PizzaToppingPlugin) SetRestaurantInformerFactory(factory informers.SharedInformerFactory) {
	p.toppingLister = factory.Restaurant().V1alpha1().Toppings().Lister()
	p.SetReadyFunc(factory.Restaurant().V1alpha1().Toppings().Informer().HasSynced)
}

func (p *PizzaToppingPlugin) ValidateInitialization() error {
	if p.toppingLister == nil {
		return fmt.Errorf("missing policy lister")
	}
	return nil
}
