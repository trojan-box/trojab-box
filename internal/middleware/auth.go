package middleware

import (
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/aresprotocols/trojan-box/internal/pkg/jwt"
	"github.com/aresprotocols/trojan-box/internal/pkg/response"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func JWTAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		response := vo.Response{Code: 0, Message: "OK"}

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Code = constant.CHECK_USER_ERROR
			response.Message = constant.MSG_CHECK_USER_ERROR
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Code = constant.CHECK_USER_ERROR
			response.Message = constant.MSG_CHECK_USER_ERROR
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		mc, err := jwt.ParseToken(parts[1], []byte(app.Conf.JwtSecret))
		if err != nil {
			response.Code = constant.CHECK_USER_ERROR
			response.Message = constant.MSG_CHECK_USER_ERROR
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Set("address", mc.Address)
		c.Next()
	}
}

func AdminAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		address := c.GetString("address")
		if address == "" {
			logger.Error("address is empty")
			response.Unauthorized(c)
			c.Abort()
			return
		}
		config := *app.Conf
		found := false
		for _, v := range config.ManagerAddress {
			if address == v {
				found = true
				break
			}
		}
		if found {
			c.Set("address", address)
			c.Next()
		} else {
			logger.Error("address not in admin address lit")
			response.NotAdmin(c)
			c.Abort()
		}
	}
}
