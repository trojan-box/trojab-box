package v1

import (
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// GetNonce godoc
// @Summary get nonce
// @Description get user nonce by address
// @Tags nonce
// @Accept json
// @Produce json
// @Param address query string true "user address"
// @Success 200 {object} vo.Response{data=string}
// @Router /nonce [get]
func GetNonce(ctx *gin.Context) {
	address, exist := ctx.GetQuery("address")
	if !exist {
		response.BadRequest(ctx)
		return
	}
	nonceUseCase := usecase.Svc.Nonce()
	nonce, err := nonceUseCase.GenerateNonce(address)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, nonce)
}
