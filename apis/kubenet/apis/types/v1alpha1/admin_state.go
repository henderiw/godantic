package v1alpha1

// +generate:validate
type AdminState string

const (
	AdminStateEnable AdminState = "enable"
	AdminStateMaintenance AdminState = "maintenance"
	AdminStateDecomissioned AdminState = "decommisioned"
	AdminStateStandby AdminState = "standby"
)