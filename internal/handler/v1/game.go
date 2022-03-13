package v1

import (
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

// @BasePath /api/v1

// PlayGame godoc
// @Summary PlayGame
// @Description play game
// @Tags game
// @Accept json
// @Produce json
// @Param . body vo.PlayGameReq false "req"
// @Success 200 {object} vo.Response{data=vo.PlayGameResp}
// @Router /game/play [post]
func PlayGame(ctx *gin.Context) {

	var playGameReq vo.PlayGameReq
	err := ctx.ShouldBind(&playGameReq)
	if err != nil {
		logger.WithError(err).Errorf("bind playGame req occur error")
		response.BadRequest(ctx)
		return
	}

	authUseCase := usecase.Svc.Auth()
	result, err := authUseCase.VerifyPlayGameRequest(playGameReq)
	if err != nil {
		logger.WithError(err).Errorf("verify playGame request occur err")
		response.BadRequestWithMsg(ctx, err.Error())
		return
	}
	if !result {
		logger.Infof("verify playGame request false")
		response.BadRequest(ctx)
		return
	}

	if app.Conf.Game.OnlyWhiteList {
		whiteList := app.Conf.WhiteList
		found := false
		for _, v := range whiteList {
			if playGameReq.Address == v {
				found = true
				break
			}
		}
		if !found {
			logger.WithError(err).Errorf("not in white list")
			response.BadRequestWithMsg(ctx, "not in white list")
			return
		}
	}

	gameUseCase := usecase.Svc.Game()
	gameSession, err := gameUseCase.PlayGame(repository.GetDB(), playGameReq, *app.Conf)
	if err != nil {
		logger.WithError(err).Errorf("play game occur err")
		response.InternalServerError(ctx, err.Error())
		return
	}
	playGameResp := vo.PlayGameResp{}
	copier.Copy(&playGameResp, gameSession)
	response.SuccessResp(ctx, playGameResp)
}

// GetGameDetail godoc
// @Summary GetGameDetail
// @Description get played game detail
// @Tags game
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param session path string true "session str"
// @Success 200 {object} vo.Response{data=vo.GameSession}
// @Router /game/{session} [get]
func GetGameDetail(ctx *gin.Context) {
	session := ctx.Param("session")
	if len(session) == 0 {
		logger.Errorf("session is emptys")
		response.BadRequest(ctx)
		return
	}
	address := ctx.GetString("address")
	if address == "" {
		logger.Error("address is empty")
		response.BadRequest(ctx)
		return
	}
	gameUseCase := usecase.Svc.Game()
	gameSession, err := gameUseCase.GetGameDetail(address, session)
	if err != nil {
		logger.WithError(err).Errorf("get game detail occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, gameSession)
}

// GetMyGameHistory godoc
// @Summary GetMyGameHistory
// @Description get my played game history
// @Tags game
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.Pagination{items=[]vo.GameHistory}}
// @Router /game/my/history [get]
func GetMyGameHistory(ctx *gin.Context) {
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
	gameUseCase := usecase.Svc.Game()
	total, histories, err := gameUseCase.GetMyHistoryByPage(address, page, size)
	if err != nil {
		logger.WithError(err).Errorf("get game history occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    histories,
	})
}

// GetGameHistories godoc
// @Summary GetGameHistory
// @Description get all user played game history
// @Tags game
// @Accept json
// @Produce json
// @Param address query string false "filter address ,if empty will query all address"
// @Param page query int false "page start from 0,default 0"
// @Param size query int false "size default 20"
// @Success 200 {object} vo.Response{data=vo.Pagination{items=[]vo.GameHistories}}
// @Router /game/histories [get]
func GetGameHistories(ctx *gin.Context) {
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
	gameUseCase := usecase.Svc.Game()
	total, histories, err := gameUseCase.GetHistoryByPageAndAddress(address, page, size)
	if err != nil {
		logger.WithError(err).Errorf("get game history occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, vo.Pagination{
		CurPage:  page,
		TotalNum: total,
		Items:    histories,
	})
}

// GetGameDetailById godoc
// @Summary GetGameDetailById
// @Description get played game detail by id
// @Tags game
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param id path string true "session id"
// @Success 200 {object} vo.Response{data=vo.GameSession}
// @Router /game/id/{id} [get]
func GetGameDetailById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if len(idStr) == 0 {
		logger.Errorf("id str is emptys")
		response.BadRequest(ctx)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.WithError(err).Errorf("Parse id str occur error")
		response.BadRequest(ctx)
		return
	}
	address := ctx.GetString("address")
	if address == "" {
		logger.Error("address is empty")
		response.BadRequest(ctx)
		return
	}
	gameUseCase := usecase.Svc.Game()
	gameSession, err := gameUseCase.GetGameDetailById(address, id)
	if err != nil {
		logger.WithError(err).Errorf("get game detail by id occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, gameSession)
}
