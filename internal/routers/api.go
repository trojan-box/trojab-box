package routers

import (
	v1 "github.com/aresprotocols/trojan-box/internal/handler/v1"
	"github.com/aresprotocols/trojan-box/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AddApi(router *gin.Engine) {

	v1g := router.Group("/api/v1")

	v1g.GET("/nonce", v1.GetNonce)

	v1gu := v1g.Group("/user")
	{
		v1gu.POST("/auth", v1.Auth)
		v1gu.GET("/profile", middleware.JWTAuth(), v1.GetUserProfile)
		v1gu.POST("/profile", middleware.JWTAuth(), v1.ModifyUserProfile)
	}
	v1gg := v1g.Group("/game")
	{
		v1gg.POST("/play", v1.PlayGame)
		v1gg.GET("/my/history", middleware.JWTAuth(), v1.GetMyGameHistory)
		v1gg.GET("/:session", middleware.JWTAuth(), v1.GetGameDetail)
		v1gg.GET("/id/:id", middleware.JWTAuth(), v1.GetGameDetailById)
		v1gg.GET("/histories", v1.GetGameHistories)
	}
	v1gb := v1g.Group("/bonus")
	{
		v1gb.GET("/my", middleware.JWTAuth(), v1.GetMyBonus)
		v1gb.GET("/my/history", middleware.JWTAuth(), v1.GetMyBonusHistory)
		v1gb.POST("/withdraw/apply", v1.WithdrawBonusApply)
		v1gb.GET("/histories", middleware.JWTAuth(), middleware.AdminAuth(), v1.GetUserBonusHistory)
	}
	v1gbc := v1g.Group("/broadcast")
	{
		v1gbc.GET("", middleware.Lang, v1.GetBroadcasts)
		v1gbc.GET("/latest", middleware.Lang, v1.GetLatestBroadcast)
	}
	v1gl := v1g.Group("/leaderboard")
	{
		v1gl.GET("", v1.GetLeaderboard)
	}

	v1gs := v1g.Group("/stats")
	{
		v1gs.GET("/daily", v1.GetDailyStats)
		v1gs.GET("/daily/list", v1.GetDailyStatsList)
		v1gs.GET("/total", v1.GetTotalStats)
		v1gs.GET("/yield/hourly", v1.GetUserYieldHourly)
	}
	v1gw := v1g.Group("/withdraw")
	{
		v1gw.GET("/histories", v1.GetWithdraws)
		v1gw.GET("/my/history", middleware.JWTAuth(), v1.GetMyWithdraws)
		v1gw.POST("/process", middleware.JWTAuth(), middleware.AdminAuth(), v1.WithdrawBonusProcess)
		v1gw.POST("/report_tx", middleware.JWTAuth(), middleware.AdminAuth(), v1.WithdrawReportHash)
	}
	v1gp := v1g.Group("/bonus_pool")
	{
		v1gp.GET("info", v1.GetBonusPoolInfo)
	}
	v1ga := v1g.Group("/app")
	{
		v1ga.GET("config", v1.GetAppConfig)
	}
	v1gss := v1g.Group("/share")
	{
		v1gss.POST("", middleware.JWTAuth(), v1.CreateSocialShare)
		v1gss.GET("", v1.GetSocialShares)
		v1gss.POST("/process", middleware.JWTAuth(), middleware.AdminAuth(), v1.SocialShareProcess)
		v1gss.GET("/my", middleware.JWTAuth(), v1.GetMySocialShares)
	}
	v1gm := v1g.Group("/message")
	{
		v1gm.GET("/my", middleware.Lang, middleware.JWTAuth(), v1.GetMyMessages)
		v1gm.POST("/read", middleware.JWTAuth(), v1.ReadMessage)
	}
	v1gf := v1g.Group("/file")
	{
		v1gf.POST("/upload", middleware.JWTAuth(), v1.UploadFile)
	}
	v1ggas := v1g.Group("/gas")
	{
		v1ggas.GET("/cal", v1.CalGasFee)
	}
}
