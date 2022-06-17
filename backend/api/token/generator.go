package token

import (
	"github.com/filedrive-team/filfind/backend/utils/jwttoken"
	logger "github.com/sirupsen/logrus"
)

func NewTokenGenerator(jwtSecret string) *jwttoken.TokenGenerator {
	ts, err := jwttoken.TokenSecretDecode(jwtSecret)
	if err != nil {
		logger.WithError(err).Fatal("call TokenSecretDecode failed")
	}
	return jwttoken.NewTokenGenerator(ts)
}

type Extend struct {
	Type string
}

type Verify func(token string, payload *jwttoken.JwtPayload) (valid bool, err error)
