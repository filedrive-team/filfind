package server

import (
	"github.com/filedrive-team/filfind/backend/api"
	"github.com/filedrive-team/filfind/backend/api/errormsg"
	"github.com/filedrive-team/filfind/backend/models"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	logger "github.com/sirupsen/logrus"
	"strings"
)

// providerList godoc
// @Summary Storage provider list
// @Tags public
// @Accept  json
// @Produce  json
// @Param page query int false "Which page to return" mininum(1)
// @Param page_size query int false "Max records to return" mininum(1) maximum(100)
// @Param sort_by query string false "Sorting option. Example: sort_by=reputation_score" Enums(reputation_score, review_score, storage_success_rate, retrieval_success_rate, price, verified_price, iso_code, quality_adj_power, storage_deals ) default(reputation_score)
// @Param order query string false "Option to order providers. Example: order=desc" Enums(asc, desc) default(desc)
// @Param sps_status query string false "Option to filter providers by registered.  Example: sps_status=all" Enums(all, registered, autofilled) default(all)
// @Param region query string false "Option to filter providers by region.  Example: region=Asia" Enums(all, Asia, Europe, Africa, Oceania, South America, North America) default(all)
// @Param raw_power_range query object false "Option to filter providers by raw_power_range(TiB) passed as a string representation of a JSON object; when implementing, make sure the parameter is URL-encoded to ensure safe transport.   Example: raw_power_range={'min':'0','max':'1024'}"
// @Param storage_success_rate_range query object false "Option to filter providers by storage_success_rate_range passed as a string representation of a JSON object; when implementing, make sure the parameter is URL-encoded to ensure safe transport.   Example: storage_success_rate_range={'min':'0.85','max':'1'}"
// @Param reputation_score_range query object false "Option to filter providers by reputation_score_range passed as a string representation of a JSON object; when implementing, make sure the parameter is URL-encoded to ensure safe transport.   Example: reputation_score_range={'min':'90','max':'100'}"
// @Param review_score_range query object false "Option to filter providers by review_score_range passed as a string representation of a JSON object; when implementing, make sure the parameter is URL-encoded to ensure safe transport.   Example: review_score_range={'min':'4','max':'5'}"
// @Param search query string false "Search providers by keyword, support miner id/name/location. Example: f01234"
// @Success 200 {array} repo.Provider
// @Router /providers [get]
func (s *Server) providerList(c *gin.Context) {
	params := new(ProviderListParam)
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

	total, list, err := s.repo.ProviderList(params.PaginationParams, params.OrderParam, params.Search, params.FilterParam)
	if err != nil {
		logger.WithError(err).Error("call ProviderList failed")
		api.JSONInternalError(c)
		return
	}

	api.JSONOk(c, &types.CommonList{
		Total: total,
		List:  list,
	})
}

// modifyProviderDetail godoc
// @Summary SP owner modify SP detailed information
// @Tags provider
// @Accept  json
// @Produce  json
// @Param Authorization header string true "jwt access token" default(Bearer YOUR_JWT)
// @Param object body ProviderDetailParams true "provider detail param"
// @Success 200
// @Router /provider/detail [post]
func (s *Server) modifyProviderDetail(c *gin.Context) {
	params := new(ProviderDetailParams)
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
	// Does user own the provider?
	owned, err := s.repo.HasProvider(tk.Uid, params.Address)
	if err != nil {
		logger.WithError(err).Error("call HasProvider failed")
		api.JSONInternalError(c)
		return
	}
	if !owned {
		api.JSONError(c, errormsg.ByCtx(c, errormsg.OwnProviderError))
		return
	}
	info := &models.ProviderInfo{
		Uid:             uid,
		Address:         params.Address,
		AvailableDeals:  params.AvailableDeals,
		Bandwidth:       params.Bandwidth,
		SealingSpeed:    params.SealingSpeed,
		ParallelDeals:   params.ParallelDeals,
		RenewableEnergy: params.RenewableEnergy,
		Certification:   params.Certification,
		IsMember:        params.IsMember,
		Experience:      params.Experience,
	}
	if err = s.repo.UpsertProviderInfo(info); err != nil {
		logger.WithError(err).Error("call UpsertProviderInfo failed")
		api.JSONInternalError(c)
		return
	}
	api.JSONOk(c, nil)
}
