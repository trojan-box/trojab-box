package v1

import (
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// GetAppConfig godoc
// @Summary GetAppConfig
// @Description get app config info
// @Tags app
// @Accept json
// @Produce json
// @Success 200 {object} vo.Response{data=vo.Config}
// @Router /app/config [get]
func GetAppConfig(ctx *gin.Context) {
	appConfig := *app.Conf
	appConfigVo := vo.Config{}
	copier.Copy(&appConfigVo, &appConfig)
	response.SuccessResp(ctx, appConfigVo)
}
