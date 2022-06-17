package repo

import (
	"encoding/json"
	"github.com/filedrive-team/filfind/backend/api/token"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/filedrive-team/filfind/backend/utils/jwttoken"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	uuid "github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"
)

type LoginUser struct {
	ClientId     string    `json:"client_id"`
	RefreshToken string    `json:"refresh_token"`
	AccessToken  string    `json:"access_token"`
	Uid          uuid.UUID `json:"uid"`
	Type         string    `json:"type"`
	Address      string    `json:"address"`
	AddressId    string    `json:"address_id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
	Logo         string    `json:"logo"`
	Location     string    `json:"location"`
}

func (m *Manager) GenerateLoginToken(user *models.User) *LoginUser {
	clientId := utils.GenerateClientId()
	ext := &token.Extend{
		Type: user.Type,
	}
	data, err := json.Marshal(ext)
	if err != nil {
		logger.Fatalf("json marshal failed,error=%v", err)
	}
	refreshToken, err := m.tkGen.Generate(settings.ProductName, jwttoken.PTRefresh, user.Uid.String(), clientId, string(data))
	if err != nil {
		logger.Fatalf("jwttoken.Generate failed,error=%v", err)

	}
	accessToken, err := m.tkGen.Generate(settings.ProductName, jwttoken.PTAccess, user.Uid.String(), clientId, string(data))
	if err != nil {
		logger.Fatalf("jwttoken.Generate failed,error=%v", err)
	}
	u := &LoginUser{
		ClientId:     clientId,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		Uid:          user.Uid,
		Type:         user.Type,
		Address:      user.AddressRobust,
		AddressId:    user.AddressId,
		Email:        user.Email,
		Name:         user.Name,
		Avatar:       user.Avatar,
		Logo:         user.Logo,
		Location:     user.Location,
	}

	return u
}

func (m *Manager) GenerateRefreshToken(uid string, clientId string, extend interface{}) string {
	refreshToken, err := m.tkGen.Generate(settings.ProductName, jwttoken.PTRefresh, uid, clientId, extend)
	if err != nil {
		logger.Fatalf("jwttoken.Generate failed,error=%v", err)
	}
	return refreshToken
}

func (m *Manager) GetTokenVerify() token.Verify {
	return m.tkGen.Verify
}
