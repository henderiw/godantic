package v1alpha1

// +generate:validate
type BFDLinkParameters struct {
	// Disabled defines if bfd is disabled or not
	Enabled *bool `json:"enabled,omitempty"`
	// MinTx defines the desired minimal interval for sending BFD packets, in msec.
	MinTx *uint32 `json:"minTx,omitempty"`
	// MinRx defines the required minimal interval for receiving BFD packets, in msec.
	MinRx *uint32 `json:"minRx,omitempty"`
	// MinEchoRx defines the echo function timer, in msec.
	MinEchoRx *uint32 `json:"minEchoRx,omitempty"`
	// Multiplier defines the number of missed packets before the session is considered down
	Multiplier *uint32 `json:"multiplier,omitempty"`
	// TTL defines the time to live on the outgoing BFD packet
	// main=2, max=255
	TTL *uint32 `json:"ttl,omitempty"`
}