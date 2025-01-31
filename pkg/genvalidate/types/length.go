package types

import (
	"fmt"
	"strings"
)

type Length struct {
	Min     *int    `json:"min,omitempty"`
	Max     *int    `json:"max,omitempty"`
	Equal   *int    `json:"equal,omitempty"`
	Message *string `json:"message,omitempty"`
	Code    *string `json:"code,omitempty"`
}

func (r *Length) String() string {
	var sb strings.Builder
	sb.WriteString("Length(")

	// Helper function to append key-value pairs
	appendField := func(name string, value interface{}) {
		if sb.Len() > len("Length(") {
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
	if r.Equal != nil {
		appendField("equal", *r.Equal)
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
