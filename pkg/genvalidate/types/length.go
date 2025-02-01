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

func (r *Length) ExpandCode(fieldName, fieldNameCode string) string {
	var sb strings.Builder
	lengthCheck := fmt.Sprintf("len(%s)", fieldNameCode)

	// Generate validation conditions
	if r.Min != nil {
		sb.WriteString(fmt.Sprintf("if %s < %d {\n", lengthCheck, *r.Min))
		sb.WriteString(generateError(fmt.Sprintf("len %s must be > %%d", fieldName), float64(*r.Min), r.Message, r.Code))
		sb.WriteString("}\n")
	}

	if r.Max != nil {
		sb.WriteString(fmt.Sprintf("if %s > %d {\n", lengthCheck, *r.Max))
		sb.WriteString(generateError(fmt.Sprintf("len %s must be < %%d", fieldName), float64(*r.Max), r.Message, r.Code))
		sb.WriteString("}\n")
	}

	if r.Equal != nil {
		sb.WriteString(fmt.Sprintf("if %s != %d {\n", lengthCheck, *r.Equal))
		sb.WriteString(generateError(fmt.Sprintf("len %s must be = %%d", fieldName), float64(*r.Equal), r.Message, r.Code))
		sb.WriteString("}\n")
	}

	return sb.String()
}

// Helper function to generate error handling code
func generateError(message string, value float64, customMsg, code *string) string {
	var sb strings.Builder
	sb.WriteString("\terrs = errors.Join(errs, fmt.Errorf(\"")

	// If a custom message is provided, use it
	if customMsg != nil {
		sb.WriteString(*customMsg)
	} else {
		sb.WriteString(message)
	}

	// Determine correct formatting: `%d` for integers, `%f` for floats
	if value == float64(int(value)) {
		sb.WriteString(`", ` + fmt.Sprintf("%d", int(value)) + ")") // Use integer format
	} else {
		sb.WriteString(`", ` + fmt.Sprintf("%f", value) + ")") // Use float format
	}

	if code != nil {
		sb.WriteString(fmt.Sprintf(" // Error code: %s", *code))
	}

	sb.WriteString(")\n")
	return sb.String()
}
