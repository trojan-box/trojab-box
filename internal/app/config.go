package app

import "time"

var (
	// Conf global app var
	Conf *Config
)

// Config global config
// nolint
type Config struct {
	Mode           string
	URL            string
	JwtSecret      string
	JwtTimeout     int
	CSRF           bool
	Debug          bool
	ManagerAddress []string
	WhiteList      []string
	HTTP           ServerConfig
	Game           GameConfig
	Ares           AresConfig
	Cron           CronConfig
	BonusPool      BonusPoolConfig
	Share          ShareConfig
	Ipfs           IpfsConfig
}

// ServerConfig server config.
type ServerConfig struct {
	Network      string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type GameConfig struct {
	MinBalance          int
	StartHour           int
	EndHour             int
	SessionLayout       string
	MinWithdraw         int64
	MinBroadcastBonus   int64
	VerifyWalletBalance bool
	OnlyWhiteList       bool
	MaxWinSmallTimes    int64
}

type AresConfig struct {
	ApiUrl               string
	AresContractAddress  string
	AresContractDecimals int
	GasStationUrl        string
}

type CronConfig struct {
	NewStar             string
	SeasonChampion      string
	DailyStats          string
	YesterdayDailyStats string
}

type BonusPoolConfig struct {
	BonusLevelAmounts      []int64
	WinBigBonusRatio       int64
	AnnualizedFactor       float64
	DefaultBonusPoolAmount int64
	PlayGameAddAmount      int64
	WithdrawReduceAmount   int64
}

type ShareConfig struct {
	CommonShareAllowChannel []int
	WithdrawShareBonus      int64
}

type IpfsConfig struct {
	Url           string
	ProjectId     string
	ProjectSecret string
}
