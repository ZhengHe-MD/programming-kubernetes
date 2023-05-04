package fuzzer

import (
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant"
	fuzz "github.com/google/gofuzz"
	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
)

var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(s *restaurant.PizzaSpec, c fuzz.Continue) {
			c.FuzzNoCustom(s) // fuzz first without calling this function again

			// avoid empty Toppings because that is defaulted
			if len(s.Toppings) == 0 {
				s.Toppings = []restaurant.PizzaTopping{
					{"salami", 1},
					{"mozzarella", 1},
					{"tomato", 1},
				}
			}

			seen := map[string]bool{}
			for i := range s.Toppings {
				// make quantity strictly positive and of reasonable size
				s.Toppings[i].Quantity = 1 + c.Intn(10)

				// remove duplicates
				for {
					if !seen[s.Toppings[i].Name] {
						break
					}
					s.Toppings[i].Name = c.RandString()
				}
				seen[s.Toppings[i].Name] = true
			}
		},
	}
}
