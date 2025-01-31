package v1alpha1

import (

	metav1 "github.com/henderiw/godantic/apis/meta/v1"
	kubenettypesv1alpha1 "github.com/henderiw/godantic/apis/kubenet/apis/types/v1alpha1"
)

// NodeSpec defines the desired state of Node
// +generate:validate
type NodeSpec struct {
	// TBD: Do we need a name here or not ??? -> right now we assume we use the name of the resource
	// the name should be defined that is unique within the system -> k8s constraint

	// Node defines the name of the node
	// +validate(length(min = 4, max = 5))
	// +validate(range(min = 4, max = 5, exclusive_min = 100))
	Node string `json:"node"`

	// *** Static immutable below ***
	PhysicalProperties kubenettypesv1alpha1.PhysicalProperties `json:",inline"`
	// *** Static immutable above ***

	// +kubebuilder:validation:Enum=`enable`;`maintenance`;`decomissioned`;`standby`;
	AdminState kubenettypesv1alpha1.AdminState `json:"adminState"`

	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined labels
	Labels map[string]string `json:"labels,omitempty"`

	// Location defines the location information where this resource is located
	// in lon/lat coordinates
	// +optional
	Location *kubenettypesv1alpha1.Location `json:"location,omitempty"`

	// Provider defines the provider implementing this resource.
	Provider *string `json:"provider,omitempty"`

	// Version define the SW version of the node
	Version *string `json:"version,omitempty"`
}

// NodeStatus defines the observed state of Node
// +generate:validate
type NodeStatus struct {
	// ConditionedStatus provides the status of the IPClain using conditions
	// - a ready condition indicates the overall status of the resource
	metav1.ConditionedStatus `json:",inline" protobuf:"bytes,1,opt,name=conditionedStatus"`
	// System ID define the unique system id of the node
	// +optional
	SystemID *string `json:"systemID,omitempty" protobuf:"bytes,2,opt,name=systemID"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories={kuid}
// A Node represents a fundamental unit that implements compute, storage, and/or networking within your environment.
// Nodes can embody physical, virtual, or containerized entities, offering versatility in deployment options to suit
// diverse infrastructure requirements.
// Nodes are logically organized within racks and sites/regions, establishing a hierarchical structure for efficient
// resource management and organization. Additionally, Nodes are associated with nodeGroups, facilitating centralized
// management and control within defined administrative boundaries.
// Each Node is assigned a provider, representing the entity responsible for implementing the specifics of the Node.
// +generate:validate
type Node struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   NodeSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status NodeStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// NodeList contains a list of Nodes
type NodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Node `json:"items" protobuf:"bytes,2,rep,name=items"`
}