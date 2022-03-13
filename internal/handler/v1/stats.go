package v1

import (
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

// GetDailyStats godoc
// @Summary GetDailyStats
// @Description get daily stats
// @Tags stats
// @Accept json
// @Produce json
// @Param day query string true "query day"
// @Success 200 {object} vo.Response{data=vo.DailyStats}
// @Router /stats/daily [get]
func GetDailyStats(ctx *gin.Context) {
	day := ctx.DefaultQuery("day", "")
	if len(day) == 0 {
		response.BadRequest(ctx)
		return
	}

	statsUseCase := usecase.Svc.Stats()

	dailyStats, err := statsUseCase.FindDailyStatsByDay(day)
	if err != nil {
		logger.WithError(err).Errorf("get daily stats error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, dailyStats)
}

// GetDailyStatsList godoc
// @Summary GetDailyStatsList
// @Description get daily stats list
// @Tags stats
// @Accept json
// @Produce json
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.DailyStats}
// @Router /stats/daily/list [get]
func GetDailyStatsList(ctx *gin.Context) {
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
	statsUseCase := usecase.Svc.Stats()
	total, dailyStats, err := statsUseCase.FindDailyStatsByPaging(page, size)
	if err != nil {
		logger.WithError(err).Errorf("get daily stats paging error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    dailyStats,
	})
}

// GetTotalStats godoc
// @Summary GetTotalStats
// @Description get total stats,will cache 1 minute
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response{data=vo.TotalStats}
// @Router /stats/total [get]
func GetTotalStats(ctx *gin.Context) {

	statsUseCase := usecase.Svc.Stats()

	totalStats, err := statsUseCase.GetTotalStats()
	if err != nil {
		logger.WithError(err).Errorf("get total stats error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, totalStats)
}

// GetUserYieldHourly godoc
// @Summary GetDailyStatsList
// @Description get user yield hourly list
// @Tags stats
// @Accept json
// @Produce json
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.UserYieldHourlyStats}
// @Router /stats/yield/hourly [get]
func GetUserYieldHourly(ctx *gin.Context) {
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
	statsUseCase := usecase.Svc.Stats()
	total, stats, err := statsUseCase.FindUserYieldHourlyStatsByPaging(page, size)
	if err != nil {
		logger.WithError(err).Errorf("get user yield hourly stats paging error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    stats,
	})
}
