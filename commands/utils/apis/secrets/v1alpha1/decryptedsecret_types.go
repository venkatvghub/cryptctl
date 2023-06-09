package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DecryptedSecret is a holding bucket for decrypted secrets so they are in
// a format that looks like a normal Secret but can't be accidentally uploaded.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DecryptedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	KeyId string            `json:"-"`
	Data  map[string]string `json:"data,omitempty"`
}

// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DecryptedSecretList contains a list of DecryptedSecret
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DecryptedSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DecryptedSecret `json:"items"`
}

func init() {

	SchemeBuilder.Register(&DecryptedSecret{}, &DecryptedSecretList{})
}
