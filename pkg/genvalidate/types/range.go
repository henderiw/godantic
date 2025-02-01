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

func (r *Range) ExpandCode(fieldName, fieldNameCode string) string {
	var sb strings.Builder

	// Generate validation conditions
	if r.Min != nil {
		sb.WriteString(fmt.Sprintf("if %s < %f {\n", fieldNameCode, *r.Min))
		sb.WriteString(generateError(fmt.Sprintf("%s must be >= %%f", fieldName), *r.Min, r.Message, r.Code))
		sb.WriteString("}\n")
	}

	if r.Max != nil {
		sb.WriteString(fmt.Sprintf("if %s > %f {\n", fieldNameCode, *r.Max))
		sb.WriteString(generateError(fmt.Sprintf("%s must be <= %%f", fieldName), *r.Max, r.Message, r.Code))
		sb.WriteString("}\n")
	}

	if r.ExclusiveMin != nil {
		sb.WriteString(fmt.Sprintf("if %s <= %f {\n", fieldNameCode, *r.ExclusiveMin))
		sb.WriteString(generateError(fmt.Sprintf("%s must be > %%f", fieldName), *r.ExclusiveMin, r.Message, r.Code))
		sb.WriteString("}\n")
	}

	if r.ExclusiveMax != nil {
		sb.WriteString(fmt.Sprintf("if %s >= %f {\n", fieldNameCode, *r.ExclusiveMax))
		sb.WriteString(generateError(fmt.Sprintf("%s must be < %%f", fieldName), *r.ExclusiveMax, r.Message, r.Code))
		sb.WriteString("}\n")
	}

	return sb.String()
}