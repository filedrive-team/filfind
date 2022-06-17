package jwttoken

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"strings"
)

type KeyType string

const (
	KTJwtHS256 KeyType = "HS256"
	KTJwtRS256 KeyType = "RS256"
)

type TokenSecret struct {
	Type       KeyType `json:"type"`
	PrivateKey []byte  `json:"private_key"`
}

func GenerateTokenSecret(kt KeyType) (ts *TokenSecret, err error) {
	switch kt {
	case KTJwtHS256:
		sk, err := ioutil.ReadAll(io.LimitReader(rand.Reader, 32))
		if err != nil {
			return nil, err
		}

		ts = &TokenSecret{
			Type:       KTJwtHS256,
			PrivateKey: sk,
		}
	default:
		err = errors.New(fmt.Sprintf("This KeyType(%s) is not supported", kt))
	}
	return
}

func TokenSecretEncode(ts *TokenSecret) (encoded string, err error) {
	var data []byte
	if data, err = json.Marshal(ts); err != nil {
		return "", err
	}
	encoded = hex.EncodeToString(data)
	return
}

func TokenSecretDecode(encoded string) (ts *TokenSecret, err error) {
	decoded, err := hex.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		return nil, err
	}

	ts = &TokenSecret{}
	if err := json.Unmarshal(decoded, ts); err != nil {
		return nil, err
	}
	return
}

func NewToken(s *TokenSecret) *Token {
	var alg jwt.Algorithm
	switch s.Type {
	case KTJwtHS256:
		alg = jwt.NewHS256(s.PrivateKey)
	}
	return &Token{
		alg: alg,
	}
}

type Token struct {
	alg jwt.Algorithm
}

func (t *Token) GenerateJwtToken(payload interface{}) (token string, err error) {
	tokenBytes, err := jwt.Sign(payload, t.alg)
	if err != nil {
		return "", err
	}
	return string(tokenBytes), nil
}

func (t *Token) VerifyJwtToken(token string, payload interface{}) error {
	if _, err := jwt.Verify([]byte(token), t.alg, payload); err != nil {
		return err
	}
	return nil
}
