package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ExtractErrorMessages(errs validator.ValidationErrors) map[string]string {
	errorMessages := make(map[string]string)
	for _, ve := range errs {
		errorMessages[ve.Field()] = fmt.Sprintf("Failed '%s' validation", ve.Tag())
	}
	return errorMessages
}
