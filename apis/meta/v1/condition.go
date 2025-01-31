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

import (
	"sort"
	"time"
)

// A ConditionType represents a condition type for a given KRM resource
type ConditionType string

// Condition Types.
const (
	// ConditionTypeReady represents the resource ready condition
	ConditionTypeReady ConditionType = "Ready"
)

// A ConditionReason represents the reason a resource is in a condition.
type ConditionReason string

// Reasons a resource is ready or not
const (
	ConditionReasonReady   ConditionReason = "Ready"
	ConditionReasonFailed  ConditionReason = "Failed"
	ConditionReasonUnknown ConditionReason = "Unknown"
)

type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in the condition.
// "ConditionFalse" means a resource is not in the condition. "ConditionUnknown" means kubernetes
// can't decide if a resource is in the condition or not. In the future, we could add other
// intermediate conditions, e.g. ConditionDegraded.
const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// +generate:validate
type Condition struct {
	// type of condition in CamelCase or in foo.example.com/CamelCase.
	// ---
	// Many .condition.type values are consistent across resources like Available, but because arbitrary conditions can be
	// useful (see .node.status.conditions), the ability to deconflict is important.
	// The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt)
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$`
	// +kubebuilder:validation:MaxLength=316
	Type string `json:"type"`
	// status of the condition, one of True, False, Unknown.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status ConditionStatus `json:"status"`
	// observedGeneration represents the .metadata.generation that the condition was set based upon.
	// For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
	// with respect to the current state of the instance.
	// +optional
	// +kubebuilder:validation:Minimum=0
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// lastTransitionTime is the last time the condition transitioned from one status to another.
	// This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	// +validate(skip)
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	// reason contains a programmatic identifier indicating the reason for the condition's last transition.
	// Producers of specific condition types may define expected values and meanings for this field,
	// and whether the values are considered a guaranteed API.
	// The value should be a CamelCase string.
	// This field may not be empty.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=`^[A-Za-z]([A-Za-z0-9_,:]*[A-Za-z0-9_])?$`
	Reason string `json:"reason"`
	// message is a human readable message indicating details about the transition.
	// This may be an empty string.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=32768
	Message string `json:"message"`
}

// Equal returns true if the condition is identical to the supplied condition,
// ignoring the LastTransitionTime.
func (c Condition) Equal(other Condition) bool {
	return c.Type == other.Type &&
		c.Status == other.Status &&
		c.Reason == other.Reason &&
		c.Message == other.Message
}

// WithMessage returns a condition by adding the provided message to existing
// condition.
func (c Condition) WithMessage(msg string) Condition {
	c.Message = msg
	return c
}

func (c Condition) IsTrue() bool {
	return c.Status == ConditionTrue
}

// A ConditionedStatus reflects the observed status of a resource. Only
// one condition of each type may exist.
// +generate:validate
type ConditionedStatus struct {
	// Conditions of the resource.
	// +optional
	// +listType:=map
	// +listMapKey:=type
	Conditions []Condition `json:"conditions,omitempty" protobuf:"bytes,1,rep,name=conditions"`
}

// NewConditionedStatus returns a stat with the supplied conditions set.
func NewConditionedStatus(c ...Condition) *ConditionedStatus {
	r := &ConditionedStatus{}
	r.SetConditions(c...)
	return r
}

// HasCondition returns if the condition is set
func (r *ConditionedStatus) HasCondition(t ConditionType) bool {
	for _, c := range r.Conditions {
		if c.Type == string(t) {
			return true
		}
	}
	return false
}

// GetCondition returns the condition for the given ConditionKind if exists,
// otherwise returns nil
func (r *ConditionedStatus) GetCondition(t ConditionType) Condition {
	for _, c := range r.Conditions {
		if c.Type == string(t) {
			return c
		}
	}
	return Condition{Type: string(t), Status: ConditionFalse}
}

// SetConditions sets the supplied conditions, replacing any existing conditions
// of the same type. This is a no-op if all supplied conditions are identical,
// ignoring the last transition time, to those already set.
func (r *ConditionedStatus) SetConditions(c ...Condition) {
	for _, new := range c {
		exists := false
		for i, existing := range r.Conditions {
			if existing.Type != new.Type {
				continue
			}

			if existing.Equal(new) {
				exists = true
				continue
			}

			r.Conditions[i] = new
			exists = true
		}
		if !exists {
			r.Conditions = append(r.Conditions, new)
		}
	}
}

// Equal returns true if the status is identical to the supplied status,
// ignoring the LastTransitionTimes and order of statuses.
func (r *ConditionedStatus) Equal(other *ConditionedStatus) bool {
	if r == nil || other == nil {
		return r == nil && other == nil
	}

	if len(other.Conditions) != len(r.Conditions) {
		return false
	}

	sc := make([]Condition, len(r.Conditions))
	copy(sc, r.Conditions)

	oc := make([]Condition, len(other.Conditions))
	copy(oc, other.Conditions)

	// We should not have more than one condition of each type.
	sort.Slice(sc, func(i, j int) bool { return sc[i].Type < sc[j].Type })
	sort.Slice(oc, func(i, j int) bool { return oc[i].Type < oc[j].Type })

	for i := range sc {
		if !sc[i].Equal(oc[i]) {
			return false
		}
	}
	return true
}

func (r *ConditionedStatus) IsConditionTrue(t ConditionType) bool {
	c := r.GetCondition(t)
	return c.Status == ConditionTrue
}

// Ready returns a condition that indicates the resource is
// ready for use.
func Ready() Condition {
	return Condition{
		Type:               string(ConditionTypeReady),
		Status:             ConditionTrue,
		LastTransitionTime: time.Now(),
		Reason:             string(ConditionReasonReady),
	}
}

// Unknown returns a condition that indicates the resource is in an
// unknown status.
func Unknown() Condition {
	return Condition{
		Type:               string(ConditionTypeReady),
		Status:             ConditionFalse,
		LastTransitionTime: time.Now(),
		Reason:             string(ConditionReasonUnknown),
	}
}

// Failed returns a condition that indicates the resource
// failed to get reconciled.
func Failed(msg string) Condition {
	return Condition{
		Type:               string(ConditionTypeReady),
		Status:             ConditionFalse,
		LastTransitionTime: time.Now(),
		Reason:             string(ConditionReasonFailed),
		Message:            msg,
	}
}
