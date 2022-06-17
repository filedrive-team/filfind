package api

import (
	"github.com/filedrive-team/filfind/backend/api/errormsg"
	"github.com/filedrive-team/filfind/backend/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JSONOk(c *gin.Context, data interface{}) {
	if data == nil {
		data = struct{}{}
	}
	c.JSON(http.StatusOK, types.CommonResp{
		Code: types.SuccessCode,
		Msg:  "OK",
		Data: data,
	})
}

func JSONError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, types.CommonResp{
		Code: types.ErrorCode,
		Msg:  msg,
	})
}

func JSONNotFound(c *gin.Context) {
	c.JSON(http.StatusOK, types.CommonResp{
		Code: types.NotFoundCode,
		Msg:  errormsg.ByCtx(c, errormsg.NotFoundError),
	})
}

func JSONInternalError(c *gin.Context) {
	c.JSON(http.StatusOK, types.CommonResp{
		Code: types.InternalErrorCode,
		Msg:  errormsg.ByCtx(c, errormsg.InternalServerError),
	})
}

func JSONExpire(c *gin.Context) {
	c.JSON(http.StatusOK, types.CommonResp{
		Code: types.ExpireCode,
		Msg:  errormsg.ByCtx(c, errormsg.Expire),
	})
}

func JSONForbidden(c *gin.Context) {
	c.JSON(http.StatusOK, types.CommonResp{
		Code: types.ForbiddenCode,
		Msg:  errormsg.ByCtx(c, errormsg.Forbidden),
	})
}

func JSONForbiddenCustom(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, types.CommonResp{
		Code: types.ForbiddenCode,
		Msg:  msg,
	})
}
