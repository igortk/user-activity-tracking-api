package utils

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

func GenerateErrorMessage[T any](data T, validation *validator.Validate) string {
	err := validation.Struct(data)
	if err == nil {
		return ""
	}

	var errors []string
	for _, fieldErr := range err.(validator.ValidationErrors) {
		var msg string
		switch fieldErr.Tag() {
		case "required":
			msg = fieldErr.Field() + " is required"
		case "gt":
			msg = fieldErr.Field() + " must be greater than " + fieldErr.Param()
		case "oneof":
			msg = fieldErr.Field() + " must be one of: " + fieldErr.Param()
		default:
			msg = fieldErr.Field() + " is invalid"
		}
		errors = append(errors, msg)
	}

	return strings.Join(errors, ". ")
}
