package v1

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

// GetLeaderboard godoc
// @Summary GetLeaderboard
// @Description get leaderboard
// @Tags leaderboard
// @Accept json
// @Produce json
// @Param type query int true "type: 1: newStar 2:season champion"
// @Success 200 {object} vo.Response{data=[]vo.LeaderboardResp}
// @Router /leaderboard [get]
func GetLeaderboard(ctx *gin.Context) {
	typeStr := ctx.DefaultQuery("type", "1")
	leaderboardType, err := strconv.Atoi(typeStr)
	if err != nil {
		response.BadRequest(ctx)
		return
	}
	leaderboardUseCase := usecase.Svc.Leaderboard()
	leaderboards, err := leaderboardUseCase.GetLeaderboardsByType(model.LeaderboardType(leaderboardType))
	if err != nil {
		logger.WithError(err).Errorf("get leaderboards error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, leaderboards)
}
