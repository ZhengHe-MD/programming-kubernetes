package install

import (
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-crd/pkg/apis/restaurant/v1alpha1"
	"github.com/ZhengHe-MD/programming-kubernetes/pizza-crd/pkg/apis/restaurant/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func Install(scheme *runtime.Scheme) {
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	utilruntime.Must(v1beta1.AddToScheme(scheme))
}
