package middleware

import (
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/gin-gonic/gin"
)

func Lang(c *gin.Context) {
	lang := ""
	for k, v := range c.Request.Header {
		if k == "Lang" {
			lang = v[0]
		}
	}
	if lang == "" {
		lang = constant.LangZhCn
	}

	c.Set("lang", lang)
}
