package types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

type ValidatorRuleParser func(string) (ValidationRule, error)

type ValidationRule interface {
	String() string
	ExpandCode(fieldName, fieldNameCode string) string
}

func InitValidationRuleRegistry() map[string]ValidatorRuleParser {
	return map[string]ValidatorRuleParser{
		"length": func(attr string) (ValidationRule, error) {
			return parseKeyValuePairs[Length](attr)
		},
		"range": func(attr string) (ValidationRule, error) {
			return parseKeyValuePairs[Range](attr)
		},
	}
	/*
		"card": func(attr string) ValidationRule {
			return parseKeyValuePairs[Card](attr)
		},
		"contains": func(attr string) ValidationRule {
			return parseKeyValuePairs[Contains](attr)
		},
		"does_not_contain": func(attr string) ValidationRule {
			return parseKeyValuePairs[DoesNotContain](attr)
		},
		"email": func(attr string) ValidationRule {
			return parseKeyValuePairs[Email](attr)
		},
	*/

	/*
		"must_match": func(attr string) ValidationRule {
			return parseKeyValuePairs[MustMatch](attr)
		},
		"regex": func(attr string) ValidationRule {
			return parseKeyValuePairs[Regex](attr)
		},
		"custom": func(attr string) ValidationRule {
			return parseKeyValuePairs[Custom](attr)
		},
	*/
}

func parseKeyValuePairs[T any](input string) (*T, error) {
	var result T
	resultValue := reflect.ValueOf(&result).Elem()


	pairs := strings.Split(input, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		// Convert key to Title Case (e.g., "min" -> "Min")
		key = strcase.ToCamel(key)

		field := resultValue.FieldByName(key)
		if !field.IsValid() || !field.CanSet() {
			fmt.Printf("Skipping invalid field: %s\n", key)
			continue
		}

		// Handle pointer types correctly
		switch field.Type().Elem().Kind() {
		case reflect.Int:
			if intValue, err := strconv.Atoi(value); err == nil {
				newValue := int(intValue)
				field.Set(reflect.ValueOf(&newValue))
			}
		case reflect.Float32, reflect.Float64:
			if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
				newValue := float64(floatValue)
				field.Set(reflect.ValueOf(&newValue))
			}
		case reflect.String:
			newValue := strings.Trim(value, `"`)
			field.Set(reflect.ValueOf(&newValue))
		default:
			fmt.Printf("Unsupported field type: %s\n", field.Type().Elem().Kind())
		}
	}

	return &result, nil
}
