package v1

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

// @BasePath /api/v1

// Auth godoc
// @Summary auth
// @Description user login and get accessToken
// @Tags user
// @Accept json
// @Produce json
// @Param . body vo.UserAuthReq false "auth"
// @Success 200 {object} vo.Response{data=string}
// @Router /user/auth [post]
func Auth(ctx *gin.Context) {

	var authReq vo.UserAuthReq
	err := ctx.ShouldBind(&authReq)
	if err != nil {
		logger.WithError(err).Errorf("bind auth requ occur error")
		response.BadRequest(ctx)
		return
	}

	authUseCase := usecase.Svc.Auth()
	result, err := authUseCase.VerifyLoginRequest(authReq)
	if err != nil {
		logger.WithError(err).Errorf("verify login request occur err")
		response.BadRequestWithMsg(ctx, err.Error())
		return
	}
	if !result {
		logger.Infof("verify auth request false")
		response.BadRequest(ctx)
		return
	}

	userUseCase := usecase.Svc.User()
	token, err := userUseCase.Auth(repository.GetDB(), authReq)
	if err != nil {
		logger.WithError(err).Errorf("login error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, token)
	// valid timestamp is no earlier than 5 minutes from the current time
	// valid nonce is existed
	// valid signed msg
	// return access token
	// destroy all nonce by address
	// insert user

}

// GetUserProfile godoc
// @Summary getUserProfile
// @Description get user profile
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Success 200 {object} vo.Response{data=vo.UserProfile}
// @Router /user/profile [get]
func GetUserProfile(ctx *gin.Context) {
	address := ctx.GetString("address")
	if address == "" {
		logger.Error("address is empty")
		response.BadRequest(ctx)
		return
	}
	userUseCase := usecase.Svc.User()
	profile, err := userUseCase.GetUserProfile(address)
	if err != nil {
		logger.WithError(err).Errorf("get user profile error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, profile)
}

// ModifyUserProfile godoc
// @Summary modifyUserProfile
// @Description modify user profile
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "accessToken"
// @Param . body vo.ModifyUserProfileReq false "modifyReq"
// @Success 200 {object} vo.Response{data=string}
// @Router /user/profile [post]
func ModifyUserProfile(ctx *gin.Context) {
	address := ctx.GetString("address")
	if address == "" {
		logger.Error("address is empty")
		response.BadRequest(ctx)
		return
	}
	req := vo.ModifyUserProfileReq{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		logger.WithError(err).Error("bind  ModifyUserProfileReq occur error")
		response.BadRequest(ctx)
		return
	}

	userUseCase := usecase.Svc.User()

	err = userUseCase.UpdateUserProfile(model.User{
		Address:  address,
		NickName: req.NickName,
		Avatar:   req.Avatar,
	})
	if err != nil {
		logger.WithError(err).Errorf("update user profile occur error")
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.SuccessResp(ctx, "")
}
