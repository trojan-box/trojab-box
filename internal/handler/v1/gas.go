package v1

import (
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

// CalGasFee godoc
// @Summary CalGasFee
// @Description calculate gas to ares
// @Tags gas
// @Accept json
// @Produce json
// @Param gas query int false "gas"
// @Success 200 {object} vo.Response{data=string}
// @Router /gas/cal [get]
func CalGasFee(ctx *gin.Context) {
	gasStr := ctx.Query("gas")
	if gasStr == "" {
		response.BadRequest(ctx)
		return
	}
	gas, err := strconv.ParseInt(gasStr, 10, 64)
	if err != nil {
		logger.WithError(err).Errorln("parse to int occur err")
		response.BadRequestWithMsg(ctx, err.Error())
		return
	}
	gasUseCase := usecase.Svc.Gas()
	aresGasFee := gasUseCase.CalGasFeeToAres(gas, *app.Conf)
	response.SuccessResp(ctx, KeepValidDecimals(aresGasFee, 7))
}

func KeepValidDecimals(num decimal.Decimal, keep int) float64 {
	if num.IsZero() {
		return 0
	}
	if num.GreaterThanOrEqual(decimal.NewFromInt(1)) {
		floatNum, _ := strconv.ParseFloat(num.StringFixed(int32(keep)), 64)
		return floatNum
	} else {
		bigNum := num.Mul(decimal.New(1, 18))
		length := len(bigNum.String())
		bigNum = bigNum.Round(0 - int32(length-keep))
		bigNum = bigNum.Div(decimal.New(1, 18))
		floatNum, _ := strconv.ParseFloat(bigNum.String(), 64)
		return floatNum
	}
}
