// GENERATED CODE - DO NOT EDIT
package v1alpha1

import (
	"errors"
	"fmt"
)

func (r ISISLevel) Validate() error {
	valid := map[string]struct{}{"L1": {}, "L2": {}, "L1L2": {}}
	if _, ok := valid[string(r)]; !ok {
		return fmt.Errorf("invalid value for ISISLevel: %s", r)
	}
	return nil
}
func (r *ISISLinkParameters) Validate() error {
	var errs error
	if r.Level != nil {
		if err := r.Level.Validate(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}
