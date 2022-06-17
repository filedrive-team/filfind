package validator

import (
	"github.com/go-playground/validator/v10"
)

func Password(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	hasDigit := false
	hasUpper := false
	hasLower := false
	hasInvalid := false
	length := len(pwd)
	for _, value := range pwd {
		switch {
		case value >= '0' && value <= '9':
			hasDigit = true
		case value >= 'A' && value <= 'Z':
			hasUpper = true
		case value >= 'a' && value <= 'z':
			hasLower = true
		default:
			hasInvalid = true
		}
	}
	return hasDigit && hasUpper && hasLower && !hasInvalid && (length >= 8 && length <= 20)
}
