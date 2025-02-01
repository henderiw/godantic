// GENERATED CODE - DO NOT EDIT
package v1alpha1

import (
	"errors"
	"fmt"
)

func (r *NodeSpec) Validate() error {
	var errs error
	if r.Node != nil {
		if len(*r.Node) < 10 {
			errs = errors.Join(errs, fmt.Errorf("len *Node must be > %d", 10))
		}
	}
	if err := r.PhysicalProperties.Validate(); err != nil {
		errs = errors.Join(errs, err)
	}
	if err := r.AdminState.Validate(); err != nil {
		errs = errors.Join(errs, err)
	}
	if r.Location != nil {
		if err := r.Location.Validate(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}
func (r *NodeStatus) Validate() error {
	return nil
}
func (r *Node) Validate() error {
	return nil
}
