package server

import (
	"context"
	"encoding/hex"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filedrive-team/filfind/backend/api"
	"github.com/filedrive-team/filfind/backend/api/errormsg"
	"github.com/filedrive-team/filfind/backend/filclient"
	"github.com/filedrive-team/filfind/backend/kvcache/vcodelimit"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/repo"
	"github.com/filedrive-team/filfind/backend/settings"
	"github.com/filedrive-team/filfind/backend/smtp"
	"github.com/filedrive-team/filfind/backend/utils/utils"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
)

// userLogin godoc
// @Summary User login
// @Tags public
// @Accept  json
// @Produce  json
// @Param object body UserParam true "user param"
// @Success 200 {object} repo.LoginUser
// @Router /userLogin [post]
func (s *Server) userLogin(c *gin.Context) {
	params := new(UserParam)
	err := c.ShouldBindJSON(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "Email") {
			code = errormsg.EmailFormatError
		} else if strings.Contains(errStr, "Password") {
			code = errormsg.PasswordFormatError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		return
	}
	user, err := s.repo.QueryUserByEmail(params.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.JSONError(c, errormsg.ByCtx(c, errormsg.Unregister))
			return
		}
		logger.WithError(err).Error("call QueryUserByEmail failed")
		api.JSONInternalError(c)
		return
	}
	// TODO: If the password is entered incorrectly for several consecutive times,
	// the password will be locked for a period of time
	if !utils.ComparePassword(user.HashedPassword, params.Password+s.conf.App.PasswordSalt) {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.WrongPassword))
		return
	}

	userInfo := s.repo.GenerateLoginToken(user)
	err = s.createLoginInfo(c, userInfo.Uid, false)
	if err != nil {
		logger.WithError(err).Error("call CreateLoginInfo failed")
	}

	api.JSONOk(c, userInfo)
}

// UserSignUp godoc
// @Summary User sign up
// @Tags public
// @Accept  json
// @Produce  json
// @Param object body SignUpParam true "sign up param"
// @Success 200 {object} repo.LoginUser
// @Router /userSignUp [post]
func (s *Server) userSignUp(c *gin.Context) {
	params := new(SignUpParam)
	err := c.ShouldBindJSON(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "Email") {
			code = errormsg.EmailFormatError
		} else if strings.Contains(errStr, "Password") {
			code = errormsg.PasswordFormatError
		} else if strings.Contains(errStr, "Type") {
			code = errormsg.UserTypeError
		} else if strings.Contains(errStr, "Address") {
			code = errormsg.AddressError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	addr, _ := address.NewFromString(params.Address)
	// verify Signature
	verifyErrCh := make(chan error)
	go func() {
		valid, err := s.verifySignature(ctx, addr, params.Message, params.Signature)
		if err == nil && !valid {
			err = errors.New("signature invalid")
		}
		if err != nil {
			logger.WithError(err).
				WithField("address", addr).
				WithField("message", params.Message).
				WithField("signature", params.Signature).
				Error("verify signature failed")
			verifyErrCh <- err
		}
		close(verifyErrCh)
	}()
	// check Address param, get address id and address robust
	var addrId, addrRobust string
	if strings.HasPrefix(params.Address, "f0") {
		robust, err := s.filClient.StateAccountKey(ctx, addr, filclient.EmptyTSK)
		if err != nil {
			logger.WithError(err).Error("call StateAccountKey failed")
			if strings.Contains(err.Error(), "unknown actor code") {
				api.JSONError(c, errormsg.ByCtx(c, errormsg.ActorNotFound))
			} else {
				api.JSONInternalError(c)
			}
			return
		}
		addrId = addr.String()
		addrRobust = robust.String()
	} else {
		id, err := s.filClient.StateLookupID(ctx, addr, filclient.EmptyTSK)
		if err != nil {
			logger.WithError(err).Warn("call StateLookupID failed")
			// debug mod ignore error
			if !s.conf.App.Debug {
				if strings.Contains(err.Error(), "actor not found") {
					api.JSONError(c, errormsg.ByCtx(c, errormsg.ActorNotFound))
				} else {
					api.JSONInternalError(c)
				}
				return
			}
		}
		addrId = id.String()
		addrRobust = addr.String()
	}

	// process result
	err = <-verifyErrCh
	if err != nil {
		logger.WithError(err).Error("call verify signature failed")
		api.JSONError(c, errormsg.ByCtx(c, errormsg.SignatureError))
		return
	}

	exist, err := s.repo.ExistUser(params.Type, addrRobust)
	if err != nil {
		logger.WithError(err).Error("call ExistUser failed")
		api.JSONInternalError(c)
		return
	}
	if exist {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.HasRegistered))
		return
	}
	existEmail, err := s.repo.ExistUserByEmail(params.Email)
	if err != nil {
		logger.WithError(err).Error("call ExistUserByEmail failed")
		api.JSONInternalError(c)
		return
	}
	if existEmail {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.EmailRegistered))
		return
	}

	hashedPassword, err := utils.GenerateHashedPassword(params.Password + s.conf.App.PasswordSalt)
	if err != nil {
		logger.WithError(err).Error("call GenerateHashedPassword failed")
		api.JSONInternalError(c)
		return
	}

	user := &models.User{
		Type:           params.Type,
		Email:          params.Email,
		HashedPassword: hashedPassword,
		Name:           params.Name,
		AddressRobust:  addrRobust,
		AddressId:      addrId,
	}
	err = s.repo.CreateUser(user)
	if err != nil {
		logger.WithError(err).Error("call CreateUser failed")
		api.JSONInternalError(c)
		return
	}

	userInfo := s.repo.GenerateLoginToken(user)
	err = s.createLoginInfo(c, userInfo.Uid, false)
	if err != nil {
		logger.WithError(err).Error("call CreateLoginInfo failed")
	}

	// async get client data cap
	if params.Type == models.ClientRole {
		go func() {
			ctxDataCap, cancelDataCap := context.WithTimeout(context.TODO(), 5*time.Minute)
			defer cancelDataCap()
			datacap := decimal.Zero
			vcs, err := s.filClient.StateVerifiedClientStatus(ctxDataCap, addr, filclient.EmptyTSK)
			if err != nil {
				logger.WithError(err).Warn("call StateVerifiedClientStatus failed")
			} else {
				datacap = decimal.NewFromBigInt(vcs.Int, 0)
			}
			cli := &models.ClientInfo{
				Uid:       userInfo.Uid,
				AddressId: addrId,
				DataCap:   datacap,
			}
			err = s.repo.UpsertClientInfoDataCap([]*models.ClientInfo{cli})
			if err != nil {
				logger.WithError(err).Error("call UpsertClientInfoDataCap failed")
			}
		}()
	}

	// send welcome message
	go func() {
		_ = s.hub.DeliverSystemMessage(userInfo.Uid.String(), settings.InitSystemMessageNotify)
	}()

	api.JSONOk(c, userInfo)
}

func (s *Server) verifySignature(ctx context.Context, addr address.Address, msgHex string, sigHex string) (valid bool, err error) {
	msg, err := hex.DecodeString(msgHex)
	if err != nil {
		return
	}
	msgStr := string(msg)
	if !strings.Contains(strings.ToLower(msgStr), strings.ToLower(settings.ProductName)) {
		err = errors.New("incorrect message")
		return
	}
	// check address in message, but ignore debug mod
	if !s.conf.App.Debug {
		if !strings.Contains(msgStr, addr.String()) {
			err = errors.New("incorrect message")
			return
		}
	}
	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return
	}

	var sig crypto.Signature
	if err = sig.UnmarshalBinary(sigBytes); err != nil {
		return
	}

	return s.filClient.WalletVerify(ctx, addr, msg, &sig)
}

func (s *Server) createLoginInfo(c *gin.Context, uid uuid.UUID, renewal bool) error {
	ua := user_agent.New(c.GetHeader("User-Agent"))
	browser, browserVersion := ua.Browser()
	loginInfo := &models.LoginInfo{
		Uid:            uid,
		Ip:             utils.GetClientIP(c),
		Browser:        browser,
		BrowserVersion: browserVersion,
		OS:             ua.OS(),
		Platform:       ua.Platform(),
		Mobile:         ua.Mobile(),
		Bot:            ua.Bot(),
		Renewal:        renewal,
	}
	return s.repo.CreateLoginInfo(loginInfo)
}

// userResetPwd godoc
// @Summary User reset password
// @Tags public
// @Accept  json
// @Produce  json
// @Param object body ResetPwdParams true "reset password param"
// @Success 200
// @Router /userResetPwd [post]
func (s *Server) userResetPwd(c *gin.Context) {
	params := new(ResetPwdParams)
	err := c.ShouldBindJSON(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "User") {
			code = errormsg.EmailFormatError
		} else if strings.Contains(errStr, "Password") {
			code = errormsg.PasswordFormatError
		} else if strings.Contains(errStr, "Vcode") {
			code = errormsg.VCodeError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		return
	}
	exist, err := s.repo.ExistUserByEmail(params.Email)
	if err != nil {
		logger.WithError(err).Error("call ExistUser failed")
		api.JSONInternalError(c)
		return
	}
	if !exist {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.Unregister))
		return
	}

	// 验证码校验
	vcode, ok := repo.GetEmailVcodeToForgetPwd(params.Email)
	if !ok {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.VCodeNotExist))
		return
	}
	if repo.IsExpireVcode(vcode) {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.VCodeExpire))
		return
	}
	if params.VCode != vcode.Code {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.VCodeError))
		return
	}

	hashedPassword, err := utils.GenerateHashedPassword(params.NewPassword + s.conf.App.PasswordSalt)
	if err != nil {
		logger.WithError(err).Error("call GenerateHashedPassword failed")
		api.JSONInternalError(c)
		return
	}
	user, err := s.repo.QueryUserByEmail(params.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.JSONError(c, errormsg.ByCtx(c, errormsg.Unregister))
			return
		}
		logger.WithError(err).Error("call QueryUserByEmail failed")
		api.JSONInternalError(c)
		return
	}
	user.HashedPassword = hashedPassword
	err = s.repo.UpdateUserPassword(user)
	if err != nil {
		logger.WithError(err).Error("call UpdateUserPassword failed")
		api.JSONInternalError(c)
		return
	}

	repo.DeleteEmailVcodeToForgetPwd(user.Email)
	api.JSONOk(c, nil)
}

// vcodeByEmailToResetPwd godoc
// @Summary Send a verification code by email to reset the password
// @Tags public
// @Accept  json
// @Produce  json
// @Param object body EmailVcodeParam true "user email param"
// @Success 200
// @Router /vcodeByEmailToResetPwd [post]
func (s *Server) vcodeByEmailToResetPwd(c *gin.Context) {
	params := new(EmailVcodeParam)
	err := c.ShouldBindJSON(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "Email") {
			code = errormsg.EmailFormatError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		return
	}
	exist, err := s.repo.ExistUserByEmail(params.Email)
	if err != nil {
		logger.WithError(err).Error("call ExistUser failed")
		api.JSONInternalError(c)
		return
	}
	if !exist {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.Unregister))
		return
	}

	if repo.IsFrequentEmailVcodeToForgetPwd(params.Email) {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.TryAgainVCodeRequest))
		return
	}
	clientIp := utils.GetClientIP(c)
	if vcodelimit.IsOutOfEmailLimit(params.Email) ||
		vcodelimit.IsOutOfIPLimit(clientIp) ||
		!vcodelimit.IncrCount(params.Email, clientIp) {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.TooManyVCodeRequest))
		return
	}

	vcode := repo.GenerateEmailVcodeToForgetPwd(params.Email, s.conf.App.Debug)
	if !s.conf.App.Debug {
		err = smtp.SendVCodeMail(smtp.ResetPwdVCodeSubject, vcode.Code, params.Email)
		if err != nil {
			logger.WithError(err).Error("call SendMail failed")

			api.JSONError(c, errormsg.ByCtx(c, errormsg.SendVCodeFailed))
			return
		}
	}

	api.JSONOk(c, nil)
}

// modifyPassword godoc
// @Summary User modify password
// @Tags user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Param object body ModifyPasswordParam true "modify password param"
// @Success 200
// @Router /user/modifyPassword [post]
func (s *Server) modifyPassword(c *gin.Context) {
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
	user, err := s.repo.QueryUserByUid(userToken.Uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.JSONError(c, errormsg.ByCtx(c, errormsg.Unregister))
			return
		}
		logger.WithError(err).Error("call QueryUserByUid failed")
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
	err = s.repo.UpdateUserPassword(user)
	if err != nil {
		logger.WithError(err).Error("call UpdateUserPassword failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, nil)
}

// Token godoc
// @Summary Generate an access token
// @Tags user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt refresh token" default(Bearer YOUR_JWT)
// @Success 200 {object} string "access token"
// @Router /user/token [post]
func (s *Server) Token(c *gin.Context) {
	token := MustGetToken(c)
	newToken := s.repo.GenerateRefreshToken(token.Uid, token.ClientId, token.Extend)
	if uid, err := uuid.FromString(token.Uid); err == nil {
		if err = s.createLoginInfo(c, uid, true); err != nil {
			logger.WithError(err).Error("call CreateLoginInfo failed")
		}
	} else {
		logger.WithError(err).Error("call uuid.FromString failed")
	}
	api.JSONOk(c, newToken)
}

// modifyProfile godoc
// @Summary User modify profile
// @Tags user
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Param object body ProfileParam true "modify profile param"
// @Success 200
// @Router /user/profile [post]
func (s *Server) modifyProfile(c *gin.Context) {
	params := new(ProfileParam)
	err := c.ShouldBindJSON(params)
	if err != nil {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.ParamsError))
		return
	}
	tk := MustGetToken(c)
	uid, err := uuid.FromString(tk.Uid)
	if err != nil {
		api.JSONForbidden(c)
		return
	}
	u := &models.User{
		Uid:          uid,
		Name:         params.Name,
		Avatar:       params.Avatar,
		Logo:         params.Logo,
		Location:     params.Location,
		ContactEmail: params.ContactEmail,
		Slack:        params.Slack,
		Github:       params.Github,
		Twitter:      params.Twitter,
		Description:  params.Description,
	}
	if err = s.repo.UpdateProfile(u); err != nil {
		logger.WithError(err).Error("call UpsertProfile failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, nil)
}
