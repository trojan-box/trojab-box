package cron

import (
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	"github.com/robfig/cron/v3"
	logger "github.com/sirupsen/logrus"
	"time"
)

func StartCron() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.WithError(err).Errorf("load lcoation Asia/Shanghai occur error")
		return
	}
	cronConfig := app.Conf.Cron
	leaderboardUseCase := usecase.Svc.Leaderboard()
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), cron.WithLocation(location))
	_, _ = c.AddFunc(cronConfig.NewStar, func() {
		leaderboardUseCase.GenerateNewStar(repository.GetDB())
	})
	_, _ = c.AddFunc(cronConfig.SeasonChampion, func() {
		leaderboardUseCase.GenerateSeasonChampion(repository.GetDB())
	})

	statsUseCase := usecase.Svc.Stats()
	_, _ = c.AddFunc(cronConfig.DailyStats, func() {
		statsUseCase.GenerateTodayDaily()
	})
	_, _ = c.AddFunc(cronConfig.YesterdayDailyStats, func() {
		statsUseCase.GenerateYesterdayDaily()
	})

	c.Start()
}
