/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "github.com/henderiw/godantic/apis/meta/v1"
	kubenetnetworkv1alpha1 "github.com/henderiw/godantic/apis/kubenet/apis/network/v1alpha1"

)

// LinkSpec defines the desired state of Link
// +generate:validate
type LinkSpec struct {
	// +kubebuilder:storageversion

	// Endpoints define the 2 endpoint identifiers of the link
	// Can only have 2 endpoints
	// +validate(length(equal = 2))
	Endpoints []*metav1.ObjectReference `json:"endpoints" protobuf:"bytes,1,opt,name=endpoints"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	Labels map[string]string `json:"labels,omitempty"`
	// BFD defines the BFD specific parameters on the link
	// +optional
	BFD *kubenetnetworkv1alpha1.BFDLinkParameters `json:"bfd,omitempty" protobuf:"bytes,3,opt,name=bfd"`
	// OSPF defines the OSPF specific parameters on the link
	// +optional
	OSPF *kubenetnetworkv1alpha1.OSPFLinkParameters `json:"ospf,omitempty" protobuf:"bytes,4,opt,name=ospf"`
	// ISIS defines the ISIS specific parameters on the link
	// +optional
	ISIS *kubenetnetworkv1alpha1.ISISLinkParameters `json:"isis,omitempty" protobuf:"bytes,5,opt,name=isis"`
	// BGP defines the BGP specific parameters on the link
	// +optional
	BGP *kubenetnetworkv1alpha1.BGPLinkParameters `json:"bgp,omitempty" protobuf:"bytes,6,opt,name=bgp"`
}

// LinkStatus defines the observed state of Link
// +generate:validate
type LinkStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	metav1.ConditionedStatus `json:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// A link represents a physical/logical connection that enables communication and data transfer
// between 2 endpoints of a node.
// +generate:validate
type Link struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   LinkSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status LinkStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// LinkList contains a list of Links
type LinkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Link `json:"items" protobuf:"bytes,2,rep,name=items"`
}
