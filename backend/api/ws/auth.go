package ws

import (
	"encoding/json"
	"github.com/filedrive-team/filfind/backend/api/token"
	"github.com/filedrive-team/filfind/backend/utils/jwttoken"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
)

func CheckAuthority(tk string, verify token.Verify, role ...string) (uid string, err error) {
	var p jwttoken.JwtPayload
	valid, err := verify(tk, &p)
	if err != nil || !valid {
		return "", errors.New(ErrExpired)
	}
	switch p.Permission {
	case jwttoken.PTAccess:
	default:
		return "", errors.New(ErrForbidden)
	}

	// check perm
	if len(role) > 0 {
		needRole := role[0]
		ext := token.Extend{}
		err = json.Unmarshal([]byte(p.Extend.(string)), &ext)
		if err != nil {
			logger.WithField("jwt.Extend", p.Extend.(string)).Errorf("json Unmarshal failed, error=%v", err)
			return "", errors.New(ErrForbidden)
		}
		if ext.Type != needRole {
			return "", errors.New(ErrForbidden)
		}
	}

	return p.Uid, nil
}
