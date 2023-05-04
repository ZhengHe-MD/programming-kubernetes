package pizza

import (
	"context"
	"fmt"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant/validation"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

func NewStrategy(typer runtime.ObjectTyper) pizzaStrategy {
	return pizzaStrategy{typer, names.SimpleNameGenerator}
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	pizza, ok := obj.(*restaurant.Pizza)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a pizza")
	}
	return pizza.ObjectMeta.Labels, generic.ObjectMetaFieldsSet(&pizza.ObjectMeta, true), nil
}

func MatchPizza(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

type pizzaStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (ps pizzaStrategy) NamespaceScoped() bool {
	return true
}

func (ps pizzaStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {}

func (ps pizzaStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (ps pizzaStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {}

func (ps pizzaStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

func (ps pizzaStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	pizza := obj.(*restaurant.Pizza)
	return validation.ValidatePizza(pizza)
}

func (ps pizzaStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (ps pizzaStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (ps pizzaStrategy) Canonicalize(obj runtime.Object) {}

func (ps pizzaStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	objPizza := obj.(*restaurant.Pizza)
	oldPizza := old.(*restaurant.Pizza)

	return validation.ValidatePizzaUpdate(objPizza, oldPizza)
}
