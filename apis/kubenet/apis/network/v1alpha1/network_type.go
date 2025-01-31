package v1alpha1

// +generate:validate
type NetworkType string

const (
	NetworkTypeP2P       NetworkType = "pointToPoint"
	NetworkTypeBroadcast NetworkType = "broadcast"
)

// +generate:validate
type Dummy int64

const (
	DummyA Dummy = iota
	DummyB
	DummyC
)
