package topping

import (
	"context"
	"fmt"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

func NewStrategy(typer runtime.ObjectTyper) toppingStrategy {
	return toppingStrategy{typer, names.SimpleNameGenerator}
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	topping, ok := obj.(*restaurant.Topping)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Topping")
	}
	return labels.Set(topping.ObjectMeta.Labels), generic.ObjectMetaFieldsSet(&topping.ObjectMeta, true), nil
}

func MatchTopping(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

type toppingStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (ts toppingStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	return
}

func (ts toppingStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

func (ts toppingStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	return
}

func (ts toppingStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (ts toppingStrategy) NamespaceScoped() bool {
	return false
}

func (ts toppingStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (ts toppingStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (ts toppingStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (ts toppingStrategy) Canonicalize(obj runtime.Object) {
}

func (ts toppingStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
