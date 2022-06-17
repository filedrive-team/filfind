package server

import (
	"github.com/filedrive-team/filfind/backend/api"
	"github.com/filedrive-team/filfind/backend/api/errormsg"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

type AdminLoginUser struct {
	ClientId     string    `json:"client_id"`
	RefreshToken string    `json:"refresh_token"`
	AccessToken  string    `json:"access_token"`
	Uid          uuid.UUID `json:"uid"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
}

// adminUserLogin godoc
// @Summary Admin user login
// @Tags admin
// @Accept  json
// @Produce  json
// @Param object body AdminUserParam true "admin user param"
// @Success 200 {object} repo.LoginUser
// @Router /admin/userLogin [post]
func (s *Server) adminUserLogin(c *gin.Context) {
	params := new(AdminUserParam)
	err := c.ShouldBindJSON(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "Name") {
			code = errormsg.UserOrPasswordError
		} else if strings.Contains(errStr, "Password") {
			code = errormsg.PasswordFormatError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		return
	}
	user, err := s.repo.QueryAdminUserByName(params.Name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.JSONError(c, errormsg.ByCtx(c, errormsg.Unregister))
			return
		}
		logger.WithError(err).Error("call QueryAdminUserByName failed")
		api.JSONInternalError(c)
		return
	}
	// TODO: If the password is entered incorrectly for several consecutive times,
	// the password will be locked for a period of time
	if !utils.ComparePassword(user.HashedPassword, params.Password+s.conf.App.PasswordSalt) {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.WrongPassword))
		return
	}
	u := &models.User{
		Uid:         user.Uid,
		Type:        user.Type,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Description: user.Description,
	}
	userInfo := s.repo.GenerateLoginToken(u)
	adminInfo := &AdminLoginUser{
		ClientId:     userInfo.ClientId,
		RefreshToken: userInfo.RefreshToken,
		AccessToken:  userInfo.AccessToken,
		Uid:          userInfo.Uid,
		Type:         userInfo.Type,
		Name:         userInfo.Name,
		Avatar:       userInfo.Avatar,
	}
	err = s.createLoginInfo(c, userInfo.Uid, false)
	if err != nil {
		logger.WithError(err).Error("call CreateLoginInfo failed")
	}

	api.JSONOk(c, adminInfo)
}

// modifyAdminPassword godoc
// @Summary Admin user modify password
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Param object body ModifyPasswordParam true "modify password param"
// @Success 200
// @Router /admin/user/modifyPassword [post]
func (s *Server) modifyAdminPassword(c *gin.Context) {
	params := new(ModifyPasswordParam)
	err := c.ShouldBindJSON(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "NewPassword") {
			code = errormsg.NewPasswordFormatError
		} else if strings.Contains(errStr, "Password") {
			code = errormsg.PasswordFormatError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		return
	}

	userToken := MustGetToken(c)
	user, err := s.repo.QueryAdminUserByUid(userToken.Uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.JSONError(c, errormsg.ByCtx(c, errormsg.Unregister))
			return
		}
		logger.WithError(err).Error("call QueryAdminUserByUid failed")
		api.JSONInternalError(c)
		return
	}
	if !utils.ComparePassword(user.HashedPassword, params.Password+s.conf.App.PasswordSalt) {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.WrongPassword))
		return
	}

	newHashedPassword, err := utils.GenerateHashedPassword(params.NewPassword + s.conf.App.PasswordSalt)
	if err != nil {
		logger.WithError(err).Error("call GenerateHashedPassword failed")
		api.JSONInternalError(c)
		return
	}
	user.HashedPassword = newHashedPassword
	err = s.repo.UpdateAdminUserPassword(user)
	if err != nil {
		logger.WithError(err).Error("call UpdateAdminUserPassword failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, nil)
}
