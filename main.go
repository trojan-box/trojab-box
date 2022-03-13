package main

import (
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/cache"
	"github.com/aresprotocols/trojan-box/internal/cron"
	"github.com/aresprotocols/trojan-box/internal/pkg/config"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/routers"
	"github.com/aresprotocols/trojan-box/internal/thirdparty"
	"github.com/aresprotocols/trojan-box/internal/usecase"
	gocache "github.com/patrickmn/go-cache"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"strconv"
	"time"
)

var (
	cfgDir = pflag.StringP("config dir", "c", "config", "config path.")
	env    = pflag.StringP("env name", "e", "", "env var name.")
)

func main() {
	pflag.Parse()
	// init config
	config.New(*cfgDir, config.WithEnv(*env))
	app.Init()
	repository.Init()
	cache.InitBonusPoolCache(repository.DB, *app.Conf)
	usecase.Svc = usecase.New(repository.DB, cache.Cache, gocache.New(5*time.Minute, 10*time.Minute))
	cron.StartCron()
	thirdparty.InitIpfs(*app.Conf)
	router := routers.NewRouter()
	logger.Info("server running")
	router.Run(":" + strconv.Itoa(int(app.Conf.HTTP.Port)))
}
