package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppSpec struct {
	Foo *string `json:"foo,omitempty"`
}

type AppStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	Ready      bool               `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type App struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AppSpec `json:"spec"`

	Status AppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type AppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []App `json:"items"`
}

func init() {
	SchemeBuilder.Register(&App{}, &AppList{})
}
