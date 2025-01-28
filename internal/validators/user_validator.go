package validators

import (
	"shiftwave-go/internal/enum"

	"github.com/go-playground/validator/v10"
)

func ValidateUserRole(fl validator.FieldLevel) bool {
	v := fl.Field()
	_, ok := enum.ParseRole(v.String())
	return ok
}
