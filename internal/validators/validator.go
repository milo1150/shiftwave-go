package validators

import (
	"fmt"
	"shiftwave-go/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ValidateSlicePayload[T any](c echo.Context, v *validator.Validate, payload *[]T) []map[string]any {
	validateErrors := []map[string]any{}

	for index, obj := range *payload {
		if err := v.Struct(obj); err != nil {
			fieldErrors := err.(validator.ValidationErrors)
			errorMessages := utils.ExtractErrorMessages(fieldErrors)
			keyMessage := fmt.Sprintf("error at index %v", index)
			validateError := map[string]any{keyMessage: errorMessages}
			validateErrors = append(validateErrors, validateError)
		}
	}

	if len(validateErrors) > 0 {
		return validateErrors
	}

	return nil
}
