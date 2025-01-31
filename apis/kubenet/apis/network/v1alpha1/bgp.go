package v1alpha1

// +generate:validate
type BGPLinkParameters struct {
	// BFD defines if BFD is enabled for BGP on this interface
	BFD *bool `json:"bfd,omitempty"`
}