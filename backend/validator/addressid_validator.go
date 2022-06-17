package validator

import (
	"github.com/filecoin-project/go-address"
	"github.com/go-playground/validator/v10"
)

func AddressId(fl validator.FieldLevel) bool {
	addrStr := fl.Field().String()
	addr, err := address.NewFromString(addrStr)
	if err != nil {
		return false
	}
	if addr.String() != addrStr {
		return false
	}
	return addr.Protocol() == address.ID
}
