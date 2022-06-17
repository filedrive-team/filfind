package server

import (
	"github.com/filedrive-team/filfind/backend/api"
	"github.com/filedrive-team/filfind/backend/api/errormsg"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	logger "github.com/sirupsen/logrus"
	"strings"
)

type MetricsOverview struct {
	RegisteredProviders int64           `json:"registered_providers"`
	AutoFilledProviders int64           `json:"auto_filled_providers"`
	RegisteredSpRatio   decimal.Decimal `json:"registered_sp_ratio"`

	InternalContacts     int64           `json:"internal_contacts"`
	AverageAccessesDaily decimal.Decimal `json:"average_accesses_daily"`
	TotalAccesses        int64           `json:"total_accesses"`

	IncrRegisteredProviders int64           `json:"incr_registered_providers"`
	IncrAutoFilledProviders int64           `json:"incr_auto_filled_providers"`
	IncrRegisteredSpRatio   decimal.Decimal `json:"incr_registered_sp_ratio"`

	IncrInternalContacts     int64           `json:"incr_internal_contacts"`
	IncrAverageAccessesDaily decimal.Decimal `json:"incr_average_accesses_daily"`
	IncrTotalAccesses        int64           `json:"incr_total_accesses"`
}

// metricsOverview godoc
// @Summary metrics overview
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Success 200 {object} MetricsOverview
// @Router /admin/metrics/overview [get]
func (s *Server) metricsOverview(c *gin.Context) {
	res, err := s.repo.MetricsOverview()
	if err != nil {
		logger.WithError(err).Error("call MetricsOverview failed")
		api.JSONInternalError(c)
		return
	}
	incrRes, err := s.repo.IncrMetricsOverviewLast24h()
	if err != nil {
		logger.WithError(err).Error("call IncrMetricsOverviewLast24h failed")
		api.JSONInternalError(c)
		return
	}

	usage := &MetricsOverview{
		RegisteredProviders: res.RegisteredProviders,
		AutoFilledProviders: res.AutoFilledProviders,
		RegisteredSpRatio:   res.RegisteredSpRatio,

		InternalContacts:     res.InternalContacts,
		AverageAccessesDaily: res.AverageAccessesDaily,
		TotalAccesses:        res.TotalAccesses,

		IncrRegisteredProviders: incrRes.RegisteredProviders,
		IncrAutoFilledProviders: incrRes.AutoFilledProviders,
		IncrRegisteredSpRatio:   incrRes.RegisteredSpRatio,

		IncrInternalContacts:     incrRes.InternalContacts,
		IncrAverageAccessesDaily: incrRes.AverageAccessesDaily,
		IncrTotalAccesses:        incrRes.TotalAccesses,
	}
	api.JSONOk(c, usage)
}

// metricsSpOverview godoc
// @Summary metrics sp overview
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Success 200 {object} repo.MetricsSpOverview
// @Router /admin/metrics/spOverview [get]
func (s *Server) metricsSpOverview(c *gin.Context) {
	res, err := s.repo.MetricsSpOverview()
	if err != nil {
		logger.WithError(err).Error("call MetricsSpOverview failed")
		api.JSONInternalError(c)
		return
	}
	api.JSONOk(c, res)
}

// metricsClientOverview godoc
// @Summary metrics client overview
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Success 200 {object} repo.MetricsClientOverview
// @Router /admin/metrics/clientOverview [get]
func (s *Server) metricsClientOverview(c *gin.Context) {
	res, err := s.repo.MetricsClientOverview()
	if err != nil {
		logger.WithError(err).Error("call MetricsClientOverview failed")
		api.JSONInternalError(c)
		return
	}
	api.JSONOk(c, res)
}

// metricsSpToDealNewClientDetail godoc
// @Summary metrics sp to deal new client detail
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Param page query int false "Which page to return" mininum(1)
// @Param page_size query int false "Max records to return" mininum(1) maximum(100)
// @Success 200 {array} repo.MetricsSpDetailItem
// @Router /admin/metrics/spToDealNewClientDetail [get]
func (s *Server) metricsSpToDealNewClientDetail(c *gin.Context) {
	params := new(types.PaginationParams)
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

	total, list, err := s.repo.MetricsSpToDealNewClientDetail(*params)
	if err != nil {
		logger.WithError(err).Error("call MetricsSpToDealNewClientDetail failed")
		api.JSONInternalError(c)
		return
	}
	api.JSONOk(c, &types.CommonList{
		Total: total,
		List:  list,
	})
}

// metricsClientToDealSpDetail godoc
// @Summary metrics client to deal sp detail
// @Tags admin
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Param page query int false "Which page to return" mininum(1)
// @Param page_size query int false "Max records to return" mininum(1) maximum(100)
// @Success 200 {array} repo.MetricsClientDetailItem
// @Router /admin/metrics/clientToDealSpDetail [get]
func (s *Server) metricsClientToDealSpDetail(c *gin.Context) {
	params := new(types.PaginationParams)
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

	total, list, err := s.repo.MetricsClientToDealSpDetail(*params)
	if err != nil {
		logger.WithError(err).Error("call MetricsClientToDealSpDetail failed")
		api.JSONInternalError(c)
		return
	}
	api.JSONOk(c, &types.CommonList{
		Total: total,
		List:  list,
	})
}
