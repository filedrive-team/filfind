package validator

import (
	"github.com/filecoin-project/go-address"
	"github.com/go-playground/validator/v10"
)

func Address(fl validator.FieldLevel) bool {
	addrStr := fl.Field().String()
	_, err := address.NewFromString(addrStr)
	if err != nil {
		return false
	}
	return true
}
