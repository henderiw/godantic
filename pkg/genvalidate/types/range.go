package types

import (
	"fmt"
	"strings"
)

type Range struct {
	Min          *float64 `json:"min,omitempty"`
	Max          *float64 `json:"max,omitempty"`
	ExclusiveMin *float64 `json:"exclusive_min,omitempty"`
	ExclusiveMax *float64 `json:"exclusive_max,omitempty"`
	Message      *string  `json:"message,omitempty"`
	Code         *string  `json:"code,omitempty"`
}

func (r *Range) String() string {
	var sb strings.Builder
	sb.WriteString("Range(")

	// Helper function to append key-value pairs
	appendField := func(name string, value interface{}) {
		if sb.Len() > len("Range(") {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s=%v", name, value))
	}

	// Append only non-nil fields
	if r.Min != nil {
		appendField("min", *r.Min)
	}
	if r.Max != nil {
		appendField("max", *r.Max)
	}
	if r.ExclusiveMin != nil {
		appendField("exclusive_min", *r.ExclusiveMin)
	}
	if r.ExclusiveMax != nil {
		appendField("exclusive_max", *r.ExclusiveMax)
	}
	if r.Message != nil {
		appendField("message", `"`+*r.Message+`"`)
	}
	if r.Code != nil {
		appendField("code", `"`+*r.Code+`"`)
	}

	sb.WriteString(")")
	return sb.String()
}
