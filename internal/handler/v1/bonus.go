package v1

import (
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

// GetMyBonus godoc
// @Summary GetMyBonus
// @Description get my bonus
// @Tags bonus
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Success 200 {object} vo.Response{data=vo.UserBonus}
// @Router /bonus/my [get]
func GetMyBonus(ctx *gin.Context) {
	address := ctx.GetString("address")
	if address == "" {
		logger.Error("address is empty")
		response.BadRequest(ctx)
		return
	}
	bonusUseCase := usecase.Svc.Bonus()
	userBonus, err := bonusUseCase.GetUserBonus(address)
	if err != nil {
		logger.WithError(err).Errorf("get user bonus error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, userBonus)
}

// GetMyBonusHistory godoc
// @Summary GetMyBonusHistory
// @Description get my bonus history
// @Tags bonus
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.Pagination{items=[]vo.BonusHistory}}
// @Router /bonus/my/history [get]
func GetMyBonusHistory(ctx *gin.Context) {
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

	bonusUseCase := usecase.Svc.Bonus()
	total, histories, err := bonusUseCase.GetBonusHistoryByPage(address, page, size, -1)
	if err != nil {
		logger.WithError(err).Errorf("get bonus history occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    histories,
	})
}

// GetUserBonusHistory godoc
// @Summary GetUserBonusHistory
// @Description get user bonus history
// @Tags bonus
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param address query string true "address"
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Param type query int false "1：win 2：withdraw"
// @Success 200 {object} vo.Response{data=vo.Pagination{items=[]vo.BonusHistory}}
// @Router /bonus/histories [get]
func GetUserBonusHistory(ctx *gin.Context) {

	address := ctx.DefaultQuery("address", "")
	if address == "" {
		logger.Error("address is empty")
		response.BadRequest(ctx)
		return
	}
	pageStr := ctx.DefaultQuery("page", "0")
	sizeStr := ctx.DefaultQuery("size", "20")
	typeStr := ctx.DefaultQuery("type", "-1")
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
	bonusRecordType, err := strconv.Atoi(typeStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}

	bonusUseCase := usecase.Svc.Bonus()
	total, histories, err := bonusUseCase.GetBonusHistoryByPage(address, page, size, bonusRecordType)
	if err != nil {
		logger.WithError(err).Errorf("get bonus history occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    histories,
	})
}

// WithdrawBonusApply godoc
// @Summary WithdrawBonusApply
// @Description apply to withdraw bonus
// @Tags bonus
// @Accept json
// @Produce json
// @Param . body vo.WithdrawBonusApplyReq false "req"
// @Success 200 {object} vo.Response{data=string}
// @Router /bonus/withdraw/apply [post]
func WithdrawBonusApply(ctx *gin.Context) {

	var withdrawApply vo.WithdrawBonusApplyReq
	err := ctx.ShouldBind(&withdrawApply)
	if err != nil {
		logger.WithError(err).Errorf("bind withdrawApply req occur error")
		response.BadRequest(ctx)
		return
	}

	authUseCase := usecase.Svc.Auth()
	result, err := authUseCase.VerifyWithdrawBonusApplyRequest(withdrawApply)
	if err != nil {
		logger.WithError(err).Errorf("verify WithdrawBonusApply request occur err")
		response.BadRequestWithMsg(ctx, err.Error())
		return
	}
	if !result {
		logger.Infof("verify WithdrawBonusApply request false")
		response.BadRequest(ctx)
		return
	}

	bonusUseCase := usecase.Svc.Bonus()
	err = bonusUseCase.WithdrawBonusApply(repository.GetDB(), withdrawApply.Address, withdrawApply.Bonus, *app.Conf)
	if err != nil {
		logger.WithError(err).Errorf("withdraw bonus apply occur err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, "")
}
