package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return err.Error()
	}

	var messages []string
	for _, e := range validationErrs {
		field := toSnakeCase(e.Field())
		switch e.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", field))
		case "email":
			messages = append(messages, fmt.Sprintf("%s must be a valid email", field))
		case "min":
			messages = append(messages, fmt.Sprintf("%s must be at least %s characters long", field, e.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s must not exceed %s characters", field, e.Param()))
		case "uuid":
			messages = append(messages, fmt.Sprintf("%s must be a valid UUID", field))
		default:
			messages = append(messages, fmt.Sprintf("%s failed on %s validation", field, e.Tag()))
		}
	}
	return strings.Join(messages, ", ")
}

func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
