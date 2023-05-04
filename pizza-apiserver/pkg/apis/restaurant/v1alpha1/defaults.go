package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_PizzaSpec(obj *PizzaSpec) {
	if len(obj.Toppings) == 0 {
		obj.Toppings = []string{"salami", "mozzarella", "tomato"}
	}
}
