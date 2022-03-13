package response

import (
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuccessResp(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &vo.Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}
func ErrorResp(ctx *gin.Context, httpStatus int, code int, message string) {
	ctx.JSON(httpStatus, &vo.Response{
		Code:    code,
		Message: message,
	})
}

func BadRequest(ctx *gin.Context) {
	ctx.JSON(http.StatusBadRequest, &vo.Response{
		Code:    constant.PARAM_NOT_TRUE_ERROR,
		Message: constant.MSG_PARAM_NOT_TRUE,
	})
}
func BadRequestWithMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, &vo.Response{
		Code:    constant.PARAM_NOT_TRUE_ERROR,
		Message: msg,
	})
}
func InternalServerError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, &vo.Response{
		Code:    constant.ERROR,
		Message: msg,
	})
}

func Unauthorized(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, vo.Response{
		Code:    constant.CHECK_USER_ERROR,
		Message: constant.MSG_CHECK_USER_ERROR,
	})
}
func NotAdmin(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, vo.Response{
		Code:    constant.CHECK_ADMIN_ERROR,
		Message: constant.MSG_CHECK_ADMIN_ERROR,
	})
}
