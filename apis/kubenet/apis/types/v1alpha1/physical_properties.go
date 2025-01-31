package v1alpha1

import "time"

// PhysicalProperties defines the properties of a physical asset
// +generate:validate
type PhysicalProperties struct {
	SerialNumber string    `json:"serialNumber"`
	Manufacturer string    `json:"manufacturer"`
	PurshaseDate time.Time `json:"purchaseDate,omitempty"`
	Type         string    `json:"type"`
}
