package v1

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

// GetMyMessages godoc
// @Summary GetMyMessages
// @Description get my message
// @Tags message
// @Accept json
// @Produce json
// @Param lang header string true "language"
// @Param Authorization header string true "accessToken"
// @Param state query int false "state 1：unRead 2:read"
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.Pagination{items=[]vo.UserMessage}}
// @Router /message/my [get]
func GetMyMessages(ctx *gin.Context) {
	lang := ctx.GetString("lang")
	address := ctx.GetString("address")
	if address == "" {
		logger.Error("address is empty")
		response.BadRequest(ctx)
		return
	}
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
	stateStr := ctx.DefaultQuery("state", "-1")
	stateInt, err := strconv.Atoi(stateStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	msgState := model.UserMessageState(stateInt)
	userMessageUseCase := usecase.Svc.UserMessage()
	total, msgs, err := userMessageUseCase.PagingByAddressAndState(lang, address, msgState, page, size)
	if err != nil {
		logger.WithError(err).Errorf("get msgs error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    msgs,
	})
}

// ReadMessage godoc
// @Summary ReadMessage
// @Description make message read
// @Tags message
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param . body vo.ReadMessage false "req"
// @Success 200 {object} vo.Response{data=string}
// @Router /message/read [post]
func ReadMessage(ctx *gin.Context) {
	address := ctx.GetString("address")
	if address == "" {
		logger.Error("address is empty")
		response.BadRequest(ctx)
		return
	}
	req := vo.ReadMessage{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		logger.WithError(err).Error("bind  ReadMessage occur error")
		response.BadRequest(ctx)
		return
	}
	userMessageUseCase := usecase.Svc.UserMessage()
	err = userMessageUseCase.ReadMessage(address, req.ID)
	if err != nil {
		logger.WithError(err).Errorf("marker message read occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, "")
}
