// GENERATED CODE - DO NOT EDIT
package v1alpha1

import (
	"fmt"
)

func (r AdminState) Validate() error {
	valid := map[string]struct{}{"enable": {}, "maintenance": {}, "decommisioned": {}, "standby": {}}
	if _, ok := valid[string(r)]; !ok {
		return fmt.Errorf("invalid value for AdminState: %s", r)
	}
	return nil
}
