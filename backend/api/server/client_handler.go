package server

import (
	"fmt"
	"github.com/filedrive-team/filfind/backend/api"
	"github.com/filedrive-team/filfind/backend/api/errormsg"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

// clientProfile godoc
// @Summary Get client profile
// @Tags public
// @Accept  json
// @Produce  json
// @Param address_id query string false "Client address id." example(f01624861)
// @Success 200 {object} repo.ClientProfile
// @Router /clientProfile [get]
func (s *Server) clientProfile(c *gin.Context) {
	params := new(AddressIdParam)
	err := c.ShouldBindQuery(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "AddressId") {
			code = errormsg.AddressError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		logger.Error(err)
		return
	}

	profile, err := s.repo.ClientProfile(params.AddressId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.JSONNotFound(c)
			return
		}
		logger.WithError(err).Error("call ClientProfile failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, profile)
}

// clientDetail godoc
// @Summary Get client detailed info
// @Tags public
// @Accept  json
// @Produce  json
// @Param address_id query string false "Client address id." example(f01624861)
// @Success 200 {object} repo.ClientInfo
// @Router /clientDetail [get]
func (s *Server) clientDetail(c *gin.Context) {
	params := new(AddressIdParam)
	err := c.ShouldBindQuery(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "AddressId") {
			code = errormsg.AddressError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		logger.Error(err)
		return
	}

	info, err := s.repo.ClientInfo(params.AddressId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.JSONNotFound(c)
			return
		}
		logger.WithError(err).Error("call ClientInfo failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, info)
}

// modifyClientDetail godoc
// @Summary client modify detailed information
// @Tags client
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Param object body ClientDetailParams true "client detail param"
// @Success 200
// @Router /client/detail [post]
func (s *Server) modifyClientDetail(c *gin.Context) {
	params := new(ClientDetailParams)
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
	u, err := s.repo.QueryUserByUid(tk.Uid)
	if err != nil {
		logger.WithError(err).Error("call QueryUserByUid failed")
		api.JSONInternalError(c)
		return
	}
	info := &models.ClientInfo{
		Uid:                uid,
		AddressId:          u.AddressId,
		Bandwidth:          params.Bandwidth,
		MonthlyStorage:     params.MonthlyStorage,
		UseCase:            params.UseCase,
		ServiceRequirement: params.ServiceRequirement,
	}
	if err = s.repo.UpsertClientInfo(info); err != nil {
		logger.WithError(err).Error("call UpsertClientInfo failed")
		api.JSONInternalError(c)
		return
	}
	api.JSONOk(c, nil)
}

// clientHistoryDealStats godoc
// @Summary Get client statistics about history deals
// @Tags public
// @Accept  json
// @Produce  json
// @Param page query int false "Which page to return" mininum(1)
// @Param page_size query int false "Max records to return" mininum(1) maximum(100)
// @Param address_id query string false "Client address id." example(f01624861)
// @Success 200 {array} repo.ClientHistoryDealStatsItem
// @Router /clientHistoryDealStats [get]
func (s *Server) clientHistoryDealStats(c *gin.Context) {
	params := new(ClientHistoryDealStatsParams)
	err := c.ShouldBindQuery(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "Page") {
			code = errormsg.PageError
		} else if strings.Contains(errStr, "PageSize") {
			code = errormsg.PageSizeError
		} else if strings.Contains(errStr, "AddressId") {
			code = errormsg.AddressError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		logger.Error(err)
		return
	}

	total, list, err := s.repo.StatsClientHistoryDeal(params.PaginationParams, params.AddressId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.JSONNotFound(c)
			return
		}
		logger.WithError(err).Error("call StatsClientHistoryDeal failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, &types.CommonList{
		Total: total,
		List:  list,
	})
}

// clientReviews godoc
// @Summary Get client reviews
// @Tags public
// @Accept  json
// @Produce  json
// @Param page query int false "Which page to return" mininum(1)
// @Param page_size query int false "Max records to return" mininum(1) maximum(100)
// @Param address_id query string false "Client address id." example(f01624861)
// @Success 200 {array} repo.Review
// @Router /clientReviews [get]
func (s *Server) clientReviews(c *gin.Context) {
	params := new(ClientReviewsParam)
	err := c.ShouldBindQuery(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "Page") {
			code = errormsg.PageError
		} else if strings.Contains(errStr, "PageSize") {
			code = errormsg.PageSizeError
		} else if strings.Contains(errStr, "AddressId") {
			code = errormsg.AddressError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		logger.Error(err)
		return
	}

	total, list, err := s.repo.ReviewsByClient(params.PaginationParams, params.AddressId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.JSONNotFound(c)
			return
		}
		logger.WithError(err).Error("call ReviewsByClient failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, &types.CommonList{
		Total: total,
		List:  list,
	})
}

// submitReview godoc
// @Summary Submit a review
// @Tags client
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Param object body ReviewParams true "review param"
// @Success 200
// @Router /client/review [post]
func (s *Server) submitReview(c *gin.Context) {
	params := new(ReviewParams)
	err := c.ShouldBindJSON(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		field := ""
		if strings.Contains(errStr, "Provider") {
			code = errormsg.AddressError
		} else if strings.Contains(errStr, "Score") {
			code = errormsg.ScoreError
		} else if strings.Contains(errStr, "Content") {
			if params.Content == "" {
				code = errormsg.CustomFieldRequiredError
			} else {
				code = errormsg.CustomTextLengthError
			}
			field = "Content"
		}
		errMsg := errormsg.ByCtx(c, code)
		if field != "" {
			errMsg = fmt.Sprintf(errMsg, field)
		}
		api.JSONError(c, errMsg)
		logger.Error(err)
		return
	}
	tk := MustGetToken(c)
	uid, err := uuid.FromString(tk.Uid)
	if err != nil {
		api.JSONForbidden(c)
		return
	}
	user, err := s.repo.QueryUserByUid(tk.Uid)
	if err != nil {
		api.JSONInternalError(c)
		return
	}
	// check deal
	has, err := s.repo.HasDeal(user.AddressId, params.Provider)
	if err != nil {
		api.JSONInternalError(c)
		return
	}
	if !has {
		api.JSONForbiddenCustom(c, errormsg.ByCtx(c, errormsg.ReviewForbidden))
		return
	}

	review := &models.Review{
		Uid:      uid,
		Client:   user.AddressId,
		Provider: params.Provider,
		Score:    params.Score,
		Content:  params.Content,
		Title:    params.Title,
	}
	if err = s.repo.InsertReview(review); err != nil {
		logger.WithError(err).Error("call InsertReview failed")
		api.JSONInternalError(c)
		return
	}
	api.JSONOk(c, nil)
}

// clientList godoc
// @Summary Client list
// @Tags public
// @Accept  json
// @Produce  json
// @Param page query int false "Which page to return" mininum(1)
// @Param page_size query int false "Max records to return" mininum(1) maximum(100)
// @Param sort_by query string false "Sorting option. Example: sort_by=data_cap" Enums(storage_deals, storage_capacity, total_data_cap, used_data_cap, data_cap, verified_deals) default(data_cap)
// @Param order query string false "Option to order providers. Example: order=desc" Enums(asc, desc) default(desc)
// @Param search query string false "Search client by keyword, support client id/name/location. Example: f01234"
// @Success 200 {array} repo.Client
// @Router /clients [get]
func (s *Server) clientList(c *gin.Context) {
	params := new(ClientListParam)
	err := c.ShouldBindQuery(params)
	if err != nil {
		errStr := err.Error()
		code := errormsg.ParamsError
		if strings.Contains(errStr, "Page") {
			code = errormsg.PageError
		} else if strings.Contains(errStr, "PageSize") {
			code = errormsg.PageSizeError
		}
		api.JSONError(c, errormsg.ByCtx(c, code))
		logger.Error(err)
		return
	}

	total, list, err := s.repo.ClientList(params.PaginationParams, params.ClientOrderParam, params.Search)
	if err != nil {
		logger.WithError(err).Error("call ClientList failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, &types.CommonList{
		Total: total,
		List:  list,
	})
}
