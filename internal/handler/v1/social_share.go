package v1

import (
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

// GetSocialShares godoc
// @Summary GetSocialShares
// @Description get  bonus withdraw ,can filter by address and type
// @Tags share
// @Accept json
// @Produce json
// @Param address query string false "filter address,if empty will return all user"
// @Param auditor_address query string false "filter auditor_address"
// @Param state query int false "state 0：unProcess 1:processing 2:processed"
// @Param type query int false "type 1:withdraw 2:common"
// @Param channel query int false "type 1:Gate 2:Weibo 3:Twitter 4:Reddit 5:Facebook"
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.Pagination{items=[]vo.SocialShare}}
// @Router /share [get]
func GetSocialShares(ctx *gin.Context) {
	address := ctx.DefaultQuery("address", "")
	auditorAddress := ctx.DefaultQuery("auditor_address", "")
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

	typeStr := ctx.DefaultQuery("type", "-1")
	typeInt, err := strconv.Atoi(typeStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	shareType := model.SocialShareType(typeInt)

	stateStr := ctx.DefaultQuery("state", "-1")
	stateInt, err := strconv.Atoi(stateStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	shareState := model.SocialShareState(stateInt)
	channelStr := ctx.DefaultQuery("channel", "-1")
	channelInt, err := strconv.Atoi(channelStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	shareChannel := model.SocialShareChannel(channelInt)
	socialShareUseCase := usecase.Svc.SocialShare()
	total, items, err := socialShareUseCase.PagingAndFilter(shareState, address, shareType,
		shareChannel, auditorAddress, -1, -1, page, size)
	if err != nil {
		logger.WithError(err).Errorf("get shares  occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    items,
	})
}

// CreateSocialShare godoc
// @Summary CreateSocialShare
// @Description add social share
// @Tags share
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param . body vo.AddSocialShareReq false "addReq"
// @Success 200 {object} vo.Response{data=string}
// @Router /share [post]
func CreateSocialShare(ctx *gin.Context) {
	address := ctx.GetString("address")
	if address == "" {
		logger.Error("address is empty")
		response.BadRequest(ctx)
		return
	}
	req := vo.AddSocialShareReq{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		logger.WithError(err).Error("bind  AddSocialShareReq occur error")
		response.BadRequest(ctx)
		return
	}

	socialShareUseCase := usecase.Svc.SocialShare()
	_, err = socialShareUseCase.CreateShare(address, req, *app.Conf)
	if err != nil {
		logger.WithError(err).Errorf("create share occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, "")
}

// SocialShareProcess godoc
// @Summary SocialShareProcess
// @Description process social share
// @Tags share
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param . body vo.ProcessSocialShareReq false "req"
// @Success 200 {object} vo.Response{data=string}
// @Router /share/process [post]
func SocialShareProcess(ctx *gin.Context) {

	auditorAddress := ctx.GetString("address")
	if auditorAddress == "" {
		logger.Error("auditorAddress is empty")
		response.BadRequest(ctx)
		return
	}
	var socialShareReq vo.ProcessSocialShareReq
	err := ctx.ShouldBind(&socialShareReq)
	if err != nil {
		logger.WithError(err).Errorf("bind socialShareReq req occur error")
		response.BadRequest(ctx)
		return
	}
	socialShareUseCase := usecase.Svc.SocialShare()
	err = socialShareUseCase.Process(repository.GetDB(), socialShareReq, auditorAddress)
	if err != nil {
		logger.WithError(err).Errorf("process social share req bonus occur err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, "")
}

// GetMySocialShares godoc
// @Summary GetMySocialShares
// @Description get my social shares
// @Tags share
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param state query int false "state 0：unProcess 1:processing 2:processed"
// @Param type query int false "type 1:withdraw 2:common"
// @Param channel query int false "type 1:Gate 2:Weibo 3:Twitter 4:Reddit 5:Facebook"
// @Param beginTime query int false "begin time"
// @Param endTime query int false "end time"
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.Pagination{items=[]vo.SocialShare}}
// @Router /share/my [get]
func GetMySocialShares(ctx *gin.Context) {
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

	typeStr := ctx.DefaultQuery("type", "-1")
	typeInt, err := strconv.Atoi(typeStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	shareType := model.SocialShareType(typeInt)

	stateStr := ctx.DefaultQuery("state", "-1")
	stateInt, err := strconv.Atoi(stateStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	shareState := model.SocialShareState(stateInt)
	channelStr := ctx.DefaultQuery("channel", "-1")
	channelInt, err := strconv.Atoi(channelStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	shareChannel := model.SocialShareChannel(channelInt)

	beginTimeStr := ctx.DefaultQuery("beginTime", "-1")
	beginTime, err := strconv.ParseInt(beginTimeStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	endTimeStr := ctx.DefaultQuery("endTime", "-1")
	endTime, err := strconv.ParseInt(endTimeStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	socialShareUseCase := usecase.Svc.SocialShare()
	total, items, err := socialShareUseCase.PagingAndFilter(shareState, address, shareType,
		shareChannel, "", beginTime, endTime, page, size)
	if err != nil {
		logger.WithError(err).Errorf("get shares  occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    items,
	})
}
