package auth

import (
	"encoding/json"
	"github.com/filedrive-team/filfind/backend/api"
	"github.com/filedrive-team/filfind/backend/api/token"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/filedrive-team/filfind/backend/utils/jwttoken"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"strings"
)

func CheckAuthority(verify token.Verify, role ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader(settings.HTTPAuthHeader)
		if !strings.HasPrefix(tokenStr, "Bearer") {
			c.Abort()
			api.JSONForbidden(c)
			return
		}
		tk := strings.TrimSpace(tokenStr[6:])
		var p jwttoken.JwtPayload
		valid, err := verify(tk, &p)
		if err != nil || !valid {
			c.Abort()
			api.JSONExpire(c)
			return
		}
		requestToken := strings.HasPrefix(c.FullPath(), settings.RefreshTokenPath)
		switch p.Permission {
		case jwttoken.PTRefresh:
			if !requestToken {
				c.Abort()
				api.JSONForbidden(c)
				return
			}
		case jwttoken.PTAccess:
			if requestToken {
				c.Abort()
				api.JSONForbidden(c)
				return
			}
		}

		// check perm
		if len(role) > 0 {
			needRole := role[0]
			ext := token.Extend{}
			err = json.Unmarshal([]byte(p.Extend.(string)), &ext)
			if err != nil {
				logger.WithField("jwt.Extend", p.Extend.(string)).Errorf("json Unmarshal failed, error=%v", err)
				c.Abort()
				api.JSONForbidden(c)
				return
			}
			if ext.Type != needRole {
				c.Abort()
				api.JSONForbidden(c)
				return
			}
		}

		c.Set(settings.TokenKey, &p)
		c.Next()
		return
	}
}
