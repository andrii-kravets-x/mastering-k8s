package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type NewResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NewResourceSpec   `json:"spec,omitempty"`
	Status NewResourceStatus `json:"status,omitempty"`
}

type NewResourceSpec struct {
	Foo string `json:"foo,omitempty"`
}

type NewResourceStatus struct {
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
type NewResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NewResource `json:"items"`
}
