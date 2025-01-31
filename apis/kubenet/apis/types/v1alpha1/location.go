package v1alpha1

// +generate:validate
type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
