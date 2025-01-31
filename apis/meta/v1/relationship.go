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

package v1

// +k8s:openapi-gen=true
// Relationships define relationships to the resource.
/*
type RelationReference struct {
	// Relationships define relationships to the resource.
	// +optional
	// +listType:=map
	// +listMapKey:=type
	// +listMapKey:=apiVersion
	// +listMapKey:=kind
	// +listMapKey=name
	Relationships []*Relationship `json:"relationships,omitempty" protobuf:"bytes,1,rep,name=relationships"`
}
	*/

// +k8s:openapi-gen=true
// Relationship define relationship parameters.
type RelationReference struct {
	// Reference defines the reference to a resource
	RelationshipReference `json:",inline"`
	Type                  string `json:"type"`
	// UserDefinedLabels define metadata to the resource.
	// defined in the spec to distingiush metadata labels from user defined label
	Labels map[string]string `json:"labels,omitempty"`
}

type RelationshipReference struct {
	APIVersion string `json:"apiVersion" protobuf:"bytes,1,opt,name=apiVersion"`
	Kind       string `json:"kind" protobuf:"bytes,3,opt,name=kind"`
	UID        string `json:"uid" protobuf:"bytes,4,opt,name=uid"`
	Name       string `json:"name" protobuf:"bytes,5,opt,name=name"`
}
