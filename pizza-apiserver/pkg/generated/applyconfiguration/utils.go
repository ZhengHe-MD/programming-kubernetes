/*
Copyright The Kubernetes Authors.

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
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	v1alpha1 "github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant/v1alpha1"
	v1beta1 "github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant/v1beta1"
	restaurantv1alpha1 "github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/generated/applyconfiguration/restaurant/v1alpha1"
	restaurantv1beta1 "github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/generated/applyconfiguration/restaurant/v1beta1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=restaurant.programming-kubernetes.info, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("Pizza"):
		return &restaurantv1alpha1.PizzaApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PizzaSpec"):
		return &restaurantv1alpha1.PizzaSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PizzaStatus"):
		return &restaurantv1alpha1.PizzaStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Topping"):
		return &restaurantv1alpha1.ToppingApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ToppingSpec"):
		return &restaurantv1alpha1.ToppingSpecApplyConfiguration{}

		// Group=restaurant.programming-kubernetes.info, Version=v1beta1
	case v1beta1.SchemeGroupVersion.WithKind("Pizza"):
		return &restaurantv1beta1.PizzaApplyConfiguration{}
	case v1beta1.SchemeGroupVersion.WithKind("PizzaSpec"):
		return &restaurantv1beta1.PizzaSpecApplyConfiguration{}
	case v1beta1.SchemeGroupVersion.WithKind("PizzaStatus"):
		return &restaurantv1beta1.PizzaStatusApplyConfiguration{}
	case v1beta1.SchemeGroupVersion.WithKind("PizzaTopping"):
		return &restaurantv1beta1.PizzaToppingApplyConfiguration{}

	}
	return nil
}
