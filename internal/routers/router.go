package routers

import (
	"github.com/aresprotocols/trojan-box/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middleware.Cors())
	AddApi(router)
	AddSwagger(router)
	return router
}
