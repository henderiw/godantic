// GENERATED CODE - DO NOT EDIT
package v1alpha1

import (
	"fmt"
)

func (r OSPFVersion) Validate() error {
	valid := map[string]struct{}{"v2": {}, "v3": {}}
	if _, ok := valid[string(r)]; !ok {
		return fmt.Errorf("invalid value for OSPFVersion: %s", r)
	}
	return nil
}
