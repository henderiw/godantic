package v1alpha1

// +generate:validate
type OSPFVersion string

const (
	OSPFVersionV2      OSPFVersion = "v2"
	OSPFVersionV3      OSPFVersion = "v3"
)

// +generate:validate
type OSPFLinkParameters struct {
	// Generic IGP Link Parameters
	IGPLinkParameters `json:",inline"`
	// Defines the OSPF area the link is assocaited with
	Area *string `json:"area,omitempty"`
}