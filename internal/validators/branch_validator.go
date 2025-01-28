package validators

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func ValidateBranchesUuid(fl validator.FieldLevel) bool {
	branches := fl.Field()

	// Check is Slice
	if branches.Kind() != reflect.Slice {
		return false
	}

	// Should not be empty slice
	if branches.Len() == 0 {
		return false
	}

	for i := 0; i < branches.Len(); i++ {
		if _, err := uuid.Parse(branches.Index(i).String()); err != nil {
			return false
		}
	}

	return true
}
