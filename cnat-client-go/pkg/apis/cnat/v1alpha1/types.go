package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	PhasePending = "PENDING"
	PhaseRunning = "RUNNING"
	PhaseDone    = "DONE"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// At runs a command at a given schedule.
type At struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtSpec   `json:"spec,omitempty"`
	Status AtStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AtList contains a list of At
type AtList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []At `json:"items"`
}

// AtSpec defines the desired state of At
type AtSpec struct {
	// Schedule...
	Schedule string `json:"schedule,omitempty"`
	// Command...
	Command string `json:"command,omitempty"`
}

// AtStatus defines the observed state of At
type AtStatus struct {
	// Phase...
	Phase string `json:"phase,omitempty"`
}
