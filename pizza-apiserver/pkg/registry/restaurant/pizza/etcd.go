package pizza

import (
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/registry"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*registry.REST, error) {
	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                   func() runtime.Object { return &restaurant.Pizza{} },
		NewListFunc:               func() runtime.Object { return &restaurant.PizzaList{} },
		PredicateFunc:             MatchPizza,
		DefaultQualifiedResource:  restaurant.Resource("pizzas"),
		SingularQualifiedResource: restaurant.Resource("pizza"),
		CreateStrategy:            strategy,
		UpdateStrategy:            strategy,
		DeleteStrategy:            strategy,
		TableConvertor:            rest.NewDefaultTableConvertor(restaurant.Resource("pizzas")),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &registry.REST{Store: store}, nil
}
