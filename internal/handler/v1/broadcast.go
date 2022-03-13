package v1

import (
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

// GetLatestBroadcast godoc
// @Summary GetLatestBroadcast
// @Description get the latest broadcast
// @Tags broadcast
// @Accept json
// @Produce json
// @Param lang header string true "language"
// @Success 200 {object} vo.Response{data=string}
// @Router /broadcast/latest [get]
func GetLatestBroadcast(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	broadcastUseCase := usecase.Svc.Broadcast()
	broadcastMsg, err := broadcastUseCase.GetLatest(lang)
	if err != nil {
		logger.WithError(err).Errorf("get broadcast error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, broadcastMsg)
}

// GetBroadcasts godoc
// @Summary GetBroadcasts
// @Description get broadcasts
// @Tags broadcast
// @Accept json
// @Produce json
// @Param lang header string true "language"
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.Pagination{items=[]string}}
// @Router /broadcast [get]
func GetBroadcasts(ctx *gin.Context) {
	lang := ctx.GetString("lang")

	pageStr := ctx.DefaultQuery("page", "0")
	sizeStr := ctx.DefaultQuery("size", "20")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	broadcastUseCase := usecase.Svc.Broadcast()
	total, broadcastMsgs, err := broadcastUseCase.GetBroadcastsByPaging(lang, page, size)
	if err != nil {
		logger.WithError(err).Errorf("get broadcasts error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    broadcastMsgs,
	})
}
