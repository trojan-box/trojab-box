package v1

import (
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
)

// GetBonusPoolInfo godoc
// @Summary GetBonusPoolInfo
// @Description get bonus pool info
// @Tags bonus_pool
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response{data=vo.BonusPoolInfo}
// @Router /bonus_pool/info [get]
func GetBonusPoolInfo(ctx *gin.Context) {
	poolCache := usecase.Svc.PoolCache()
	bonusAmount := poolCache.GetBonusAmountFromPool()
	bonusInfo := vo.BonusPoolInfo{Total: bonusAmount}
	response.SuccessResp(ctx, bonusInfo)
}
