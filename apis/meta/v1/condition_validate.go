// GENERATED CODE - DO NOT EDIT
package v1

import (
	"errors"
	"fmt"
)

func (r ConditionType) Validate() error {
	valid := map[string]struct{}{"Ready": {}}
	if _, ok := valid[string(r)]; !ok {
		return fmt.Errorf("invalid value for ConditionType: %s", r)
	}
	return nil
}
func (r ConditionReason) Validate() error {
	valid := map[string]struct{}{"Ready": {}, "Failed": {}, "Unknown": {}}
	if _, ok := valid[string(r)]; !ok {
		return fmt.Errorf("invalid value for ConditionReason: %s", r)
	}
	return nil
}
func (r ConditionStatus) Validate() error {
	valid := map[string]struct{}{"True": {}, "False": {}, "Unknown": {}}
	if _, ok := valid[string(r)]; !ok {
		return fmt.Errorf("invalid value for ConditionStatus: %s", r)
	}
	return nil
}
func (r *Condition) Validate() error {
	var errs error
	if err := r.Status.Validate(); err != nil {
		errs = errors.Join(errs, err)
	}
	if errs != nil {
		return errs
	}
	return nil
}
func (r *ConditionedStatus) Validate() error {
	var errs error
	for _, item := range r.Conditions {
		if err := item.Validate(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}
