package validator

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

func User(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	hasInvalid := false
	if strings.Contains(str, "@") {
		// TODO check email format
		return true
	}

	length := len(str)
	for _, value := range str {
		switch {
		case value >= '0' && value <= '9':
		case value >= 'A' && value <= 'Z':
		case value >= 'a' && value <= 'z':
		default:
			hasInvalid = true
		}
	}
	return !hasInvalid && (length >= 8 && length <= 32)
}
