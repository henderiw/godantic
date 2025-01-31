// GENERATED CODE - DO NOT EDIT
package v1alpha1

import (
	"fmt"
)

func (r NetworkType) Validate() error {
	valid := map[string]struct{}{"pointToPoint": {}, "broadcast": {}}
	if _, ok := valid[string(r)]; !ok {
		return fmt.Errorf("invalid value for NetworkType: %s", r)
	}
	return nil
}
func (r Dummy) Validate() error {
	valid := map[int64]struct{}{0: {}, 1: {}, 2: {}}
	if _, ok := valid[int64(r)]; !ok {
		return fmt.Errorf("invalid value for Dummy: %d", r)
	}
	return nil
}
