package server

import (
	"github.com/filedrive-team/filfind/backend/api"
	"github.com/filedrive-team/filfind/backend/api/errormsg"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

// spOwnerProfile godoc
// @Summary Get SP owner profile
// @Tags public
// @Accept  json
// @Produce  json
// @Param address_id query string false "Owner address id." example(f01694606)
// @Success 200 {object} repo.SPOwnerProfile
// @Router /spOwnerProfile [get]
func (s *Server) spOwnerProfile(c *gin.Context) {
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

	profile, err := s.repo.SPOwnerProfile(params.AddressId)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.WithError(err).Error("call SPOwnerProfile failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, profile)
}

// spServiceDetail godoc
// @Summary Get SP service detail list by owner
// @Tags public
// @Accept  json
// @Produce  json
// @Param address_id query string false "Owner address id." example(f01694606)
// @Success 200 {array} repo.Provider
// @Router /spServiceDetail [get]
func (s *Server) spServiceDetail(c *gin.Context) {
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

	profile, err := s.repo.ProviderListByOwner(params.AddressId)
	if err != nil {
		logger.WithError(err).Error("call SPOwnerProfile failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, profile)
}

// spOwnerReviews godoc
// @Summary Get reviews by SP owner
// @Tags public
// @Accept  json
// @Produce  json
// @Param page query int false "Which page to return" mininum(1)
// @Param page_size query int false "Max records to return" mininum(1) maximum(100)
// @Param address_id query string false "Owner address id." example(f01694606)
// @Success 200 {array} repo.Review
// @Router /spOwnerReviews [get]
func (s *Server) spOwnerReviews(c *gin.Context) {
	params := new(spOwnerReviewsParam)
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

	total, list, err := s.repo.ReviewsBySpOwner(params.PaginationParams, params.AddressId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.JSONNotFound(c)
			return
		}
		logger.WithError(err).Error("call ReviewsBySpOwner failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, &types.CommonList{
		Total: total,
		List:  list,
	})
}
