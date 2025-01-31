package v1alpha1

// +generate:validate
type ISISLevel string

const (
	ISISLevelL1   ISISLevel = "L1"
	ISISLevelL2   ISISLevel = "L2"
	ISISLevelL1L2 ISISLevel = "L1L2"
)

// +generate:validate
type ISISLinkParameters struct {
	// Generic IGP Link Parameters
	IGPLinkParameters `json:",inline"`
	// Defines the ISIS level the link is assocaited with
	Level *ISISLevel `json:"area,omitempty"`
}
