// GENERATED CODE - DO NOT EDIT
package v1alpha1

import (
	"errors"
	"fmt"
)

func (r *LinkSpec) Validate() error {
	var errs error
	if len(r.Endpoints) != 2 {
		errs = errors.Join(errs, fmt.Errorf("len Endpoints must be = %d", 2))
	}
	for _, item := range r.Endpoints {
		if item != nil {
			if err := item.Validate(); err != nil {
				errs = errors.Join(errs, err)
			}
		}
	}
	if r.BFD != nil {
		if err := r.BFD.Validate(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if r.OSPF != nil {
		if err := r.OSPF.Validate(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if r.ISIS != nil {
		if err := r.ISIS.Validate(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if r.BGP != nil {
		if err := r.BGP.Validate(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}
func (r *LinkStatus) Validate() error {
	return nil
}
func (r *Link) Validate() error {
	return nil
}
