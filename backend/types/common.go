package types

import (
	"net/http"
)

// code for common response type
const (
	SuccessCode       = http.StatusOK
	ErrorCode         = http.StatusBadRequest
	InternalErrorCode = http.StatusInternalServerError
	ForbiddenCode     = http.StatusForbidden
	ExpireCode        = http.StatusUnauthorized
	NotFoundCode      = http.StatusNotFound
)

// CommonResp - api common response type
// success - {"code": 200, "msg": "", "data": {}}
// error - {"code": 400, "msg": "", "data": {}}
// forbidden - {"code": 403, "msg": "", "data": {}}
// expire - {"code": 401, "msg": "", "data": {}}
type CommonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type CommonList struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

type PaginationParams struct {
	Page     int `json:"page" form:"page" binding:"required,min=1" example:"1"`
	PageSize int `json:"page_size" form:"page_size" binding:"required,min=1,max=100" example:"20"`
}
