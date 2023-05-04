package validation

import (
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidatePizza(f *restaurant.Pizza) (errs field.ErrorList) {
	errs = append(errs, ValidatePizzaSpec(&f.Spec, field.NewPath("spec"))...)
	return errs
}

func ValidatePizzaUpdate(obj *restaurant.Pizza, old *restaurant.Pizza) (errs field.ErrorList) {
	// TODO
	return errs
}

func ValidatePizzaSpec(s *restaurant.PizzaSpec, fldPath *field.Path) (errs field.ErrorList) {
	seen := make(map[string]bool)
	for i := range s.Toppings {
		if s.Toppings[i].Quantity <= 0 {
			errs = append(errs, field.Invalid(
				fldPath.Child("toppings").Index(i).Child("quantity"),
				s.Toppings[i].Quantity,
				"cannot be negative or zero"))
		}

		if len(s.Toppings[i].Name) == 0 {
			errs = append(errs, field.Invalid(
				fldPath.Child("toppings").Index(i).Child("name"),
				s.Toppings[i].Name,
				"cannot be empty"))
		} else {
			if seen[s.Toppings[i].Name] {
				errs = append(errs, field.Invalid(
					fldPath.Child("toppings").Index(i).Child("name"),
					s.Toppings[i].Name,
					"must be unique"))
			}
			seen[s.Toppings[i].Name] = true
		}
	}

	return
}
