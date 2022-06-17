package jwttoken

import (
	"github.com/gbrlsnchs/jwt/v3"
	"sync"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	DefaultSecret = "7b2274797065223a224853323536222c22707269766174655f6b6579223a2236327254426361655a626835556e31453336724d55625267397773796e305357786d744178636b733069493d227d"

	PTAccess  PermissionType = "access"  // Access to all interfaces except the generation of access Token interface
	PTRefresh PermissionType = "refresh" // Only the generated Access Token interface can be accessed
	PTCustom  PermissionType = "custom"  // Custom, not set the ExpirationTime field

	AccessTokenDuration  = time.Hour * 1
	RefreshTokenDuration = time.Hour * 24 * 7
)

type PermissionType string

type JwtPayload struct {
	jwt.Payload
	Permission PermissionType `json:"perm"`
	Uid        string         `json:"uid"`
	ClientId   string         `json:"cid,omitempty"`
	Extend     interface{}    `json:"extend,omitempty"`
}

var once sync.Once
var generator *TokenGenerator

// GetDefaultTokenGenerator return a singleton TokenGenerator with default secret
func GetDefaultTokenGenerator() *TokenGenerator {
	once.Do(func() {
		ts, _ := TokenSecretDecode(DefaultSecret)
		generator = NewTokenGenerator(ts)
	})
	return generator
}

func NewTokenGenerator(ts *TokenSecret) *TokenGenerator {
	return &TokenGenerator{
		tk: NewToken(ts),
	}
}

type TokenGenerator struct {
	tk *Token
}

func (tg *TokenGenerator) Generate(issuer string, perm PermissionType, uid string, clientId string, extend interface{}) (token string, err error) {
	var d time.Duration
	switch perm {
	case PTAccess:
		d = AccessTokenDuration
	case PTRefresh:
		d = RefreshTokenDuration
	case PTCustom:
	}
	p := &JwtPayload{}
	p.JWTID = uuid.NewV4().String()
	p.Uid = uid
	p.ClientId = clientId
	if perm != PTCustom {
		p.ExpirationTime = jwt.NumericDate(time.Now().Add(d))
	}
	p.Permission = perm
	p.Issuer = issuer
	p.IssuedAt = jwt.NumericDate(time.Now())
	p.Extend = extend
	return tg.tk.GenerateJwtToken(p)
}

func (tg *TokenGenerator) Verify(token string, payload *JwtPayload) (valid bool, err error) {
	err = tg.tk.VerifyJwtToken(token, payload)
	if err != nil {
		return false, err
	}
	if payload.Permission != PTCustom {
		if payload.ExpirationTime == nil {
			return false, errors.New("ExpirationTime field is nil")
		}
		if payload.ExpirationTime.Time.Before(time.Now()) {
			return false, nil
		}
	}
	return true, nil
}
