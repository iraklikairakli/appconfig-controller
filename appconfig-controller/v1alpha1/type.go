package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppConfigSpec defines the desired state of AppConfig
type AppConfigSpec struct {
	Environment string            `json:"environment,omitempty"`
	Settings    map[string]string `json:"settings,omitempty"`
}

// AppConfigStatus defines the observed state of AppConfig
type AppConfigStatus struct {
	Applied bool `json:"applied,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AppConfig is the Schema for the appconfigs API
type AppConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppConfigSpec   `json:"spec,omitempty"`
	Status AppConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppConfigList contains a list of AppConfig
type AppConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppConfig `json:"items"`
}
