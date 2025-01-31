package v1

import "time"

// ManagedFieldsEntry is a workflow-id, a FieldSet and the group version of the resource
// that the fieldset applies to.
type ManagedFieldsEntry struct {
	// Manager is an identifier of the workflow managing these fields.
	Manager string `json:"manager,omitempty" protobuf:"bytes,1,opt,name=manager"`
	// Operation is the type of operation which lead to this ManagedFieldsEntry being created.
	// The only valid values for this field are 'Apply' and 'Update'.
	Operation ManagedFieldsOperationType `json:"operation,omitempty" protobuf:"bytes,2,opt,name=operation,casttype=ManagedFieldsOperationType"`
	// APIVersion defines the version of this resource that this field set
	// applies to. The format is "group/version" just like the top-level
	// APIVersion field. It is necessary to track the version of a field
	// set because it cannot be automatically converted.
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,3,opt,name=apiVersion"`
	// Time is the timestamp of when the ManagedFields entry was added. The
	// timestamp will also be updated if a field is added, the manager
	// changes any of the owned fields value or removes a field. The
	// timestamp does not update when a field is removed from the entry
	// because another manager took it over.
	// +optional
	Time *time.Time `json:"time,omitempty" protobuf:"bytes,4,opt,name=time"`

	// Fields is tombstoned to show why 5 is a reserved protobuf tag.
	//Fields *Fields `json:"fields,omitempty" protobuf:"bytes,5,opt,name=fields,casttype=Fields"`

	// FieldsType is the discriminator for the different fields format and version.
	// There is currently only one possible value: "FieldsV1"
	FieldsType string `json:"fieldsType,omitempty" protobuf:"bytes,6,opt,name=fieldsType"`
	// FieldsV1 holds the first JSON version format as described in the "FieldsV1" type.
	// +optional
	FieldsV1 *FieldsV1 `json:"fieldsV1,omitempty" protobuf:"bytes,7,opt,name=fieldsV1"`

	// Subresource is the name of the subresource used to update that object, or
	// empty string if the object was updated through the main resource. The
	// value of this field is used to distinguish between managers, even if they
	// share the same name. For example, a status update will be distinct from a
	// regular update using the same manager name.
	// Note that the APIVersion field is not related to the Subresource field and
	// it always corresponds to the version of the main resource.
	Subresource string `json:"subresource,omitempty" protobuf:"bytes,8,opt,name=subresource"`
}

// ManagedFieldsOperationType is the type of operation which lead to a ManagedFieldsEntry being created.
type ManagedFieldsOperationType string

const (
	ManagedFieldsOperationApply  ManagedFieldsOperationType = "Apply"
	ManagedFieldsOperationUpdate ManagedFieldsOperationType = "Update"
)

// FieldsV1 stores a set of fields in a data structure like a Trie, in JSON format.
//
// Each key is either a '.' representing the field itself, and will always map to an empty set,
// or a string representing a sub-field or item. The string will follow one of these four formats:
// 'f:<name>', where <name> is the name of a field in a struct, or key in a map
// 'v:<value>', where <value> is the exact json formatted value of a list item
// 'i:<index>', where <index> is position of a item in a list
// 'k:<keys>', where <keys> is a map of  a list item's key fields to their unique values
// If a key maps to an empty Fields value, the field that key represents is part of the set.
//
// The exact format is defined in sigs.k8s.io/structured-merge-diff
// +protobuf.options.(gogoproto.goproto_stringer)=false
type FieldsV1 struct {
	// Raw is the underlying serialization of this object.
	Raw []byte `json:"-" protobuf:"bytes,1,opt,name=Raw"`
}