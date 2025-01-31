package v1alpha1

// +generate:validate
type IGPLinkParameters struct {
	// Type defines the type of network
	// enum broadcast, pointToPoint;
	// default=pointToPoint
	NetworkType *NetworkType `json:"networkType,omitempty"`
	// Passive defines if this interface is passive
	Passive *bool `json:"passive,omitempty"`
	// BFD defines if BFD is enabled for the IGP on this interface
	// default:=true
	BFD *bool `json:"bfd,omitempty"`
	// Metric defines the interface metric associated with the native routing topology
	Metric *uint32 `json:"metric,omitempty"`
}