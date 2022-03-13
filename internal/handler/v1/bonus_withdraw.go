package v1

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

// GetWithdraws godoc
// @Summary GetWithdraws
// @Description get  bonus withdraw ,can filter by address and type
// @Tags withdraw
// @Accept json
// @Produce json
// @Param address query string false "filter address,if empty will return all user"
// @Param type query int false "type 0：unProcess 1:processing 2:processed"
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.Pagination{items=[]vo.BonusWithdraw}}
// @Router /withdraw/histories [get]
func GetWithdraws(ctx *gin.Context) {
	address := ctx.DefaultQuery("address", "")
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
	withdrawType := model.BonusWithdrawState(typeInt)
	bonusWithdrawUseCase := usecase.Svc.BonusWithdraw()
	total, items, err := bonusWithdrawUseCase.PagingByStateAndAddress(withdrawType, address, page, size)
	if err != nil {
		logger.WithError(err).Errorf("get bonus withdraw by type occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    items,
	})
}

// WithdrawBonusProcess godoc
// @Summary WithdrawBonusProcess
// @Description process withdraw bonus
// @Tags withdraw
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param . body vo.ProcessWithdrawBonusReq false "req"
// @Success 200 {object} vo.Response{data=string}
// @Router /withdraw/process [post]
func WithdrawBonusProcess(ctx *gin.Context) {

	var processWithdrawReq vo.ProcessWithdrawBonusReq
	err := ctx.ShouldBind(&processWithdrawReq)
	if err != nil {
		logger.WithError(err).Errorf("bind processWithdrawReq req occur error")
		response.BadRequest(ctx)
		return
	}

	bonusWithdrawUseCase := usecase.Svc.BonusWithdraw()
	err = bonusWithdrawUseCase.ProcessWithdraw(repository.GetDB(), processWithdrawReq.ID, processWithdrawReq.Txhash)
	if err != nil {
		logger.WithError(err).Errorf("process withdraw bonus occur err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, "")
}

// GetMyWithdraws godoc
// @Summary GetMyWithdraws
// @Description get my bonus withdraw ,can filter by type
// @Tags withdraw
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param type query int false "type 0：unProcess 1:processing 2:processed"
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.Pagination{items=[]vo.BonusWithdraw}}
// @Router /withdraw/my/history [get]
func GetMyWithdraws(ctx *gin.Context) {
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
	withdrawType := model.BonusWithdrawState(typeInt)
	bonusWithdrawUseCase := usecase.Svc.BonusWithdraw()
	total, items, err := bonusWithdrawUseCase.PagingByStateAndAddress(withdrawType, address, page, size)
	if err != nil {
		logger.WithError(err).Errorf("get bonus withdraw by type occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    items,
	})
}

// WithdrawReportHash godoc
// @Summary WithdrawReportHash
// @Description report processed withdraw tx hash
// @Tags withdraw
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param . body vo.ProcessWithdrawBonusReq false "req"
// @Success 200 {object} vo.Response{data=string}
// @Router /withdraw/report_tx [post]
func WithdrawReportHash(ctx *gin.Context) {

	var processWithdrawReq vo.ProcessWithdrawBonusReq
	err := ctx.ShouldBind(&processWithdrawReq)
	if err != nil {
		logger.WithError(err).Errorf("bind processWithdrawReq req occur error")
		response.BadRequest(ctx)
		return
	}

	bonusWithdrawUseCase := usecase.Svc.BonusWithdraw()
	err = bonusWithdrawUseCase.ReportProcessedWithdrawTxhash(repository.GetDB(), processWithdrawReq.ID, processWithdrawReq.Txhash)
	if err != nil {
		logger.WithError(err).Errorf("report withdraw txhash occur err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, "")
}
