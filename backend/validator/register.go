package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	logger "github.com/sirupsen/logrus"
)

func InitExtendValidation() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		logger.Fatal("validator engine changed")
		return
	}

	// Register new custom validation rules
	var err error
	for _, e := range []struct {
		tag string
		fn  func(fl validator.FieldLevel) bool
	}{
		{"user", User},
		{"password", Password},
		{"address", Address},
		{"addressid", AddressId},
	} {
		err = v.RegisterValidation(e.tag, e.fn)
		if err != nil {
			logger.WithField("tag", e.tag).WithError(err).Fatal("validator")
		}
	}
}
