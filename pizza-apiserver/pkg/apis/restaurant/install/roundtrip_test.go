package install

import (
	restaurantfuzzer "github.com/ZhengHe-MD/programming-kubernetes/pizza-apiserver/pkg/apis/restaurant/fuzzer"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"
	metafuzzer "k8s.io/apimachinery/pkg/apis/meta/fuzzer"
	"testing"
)

func TestRoundTripTypes(t *testing.T) {
	roundtrip.RoundTripTestForAPIGroup(t, Install, fuzzer.MergeFuzzerFuncs(metafuzzer.Funcs, restaurantfuzzer.Funcs))
}
