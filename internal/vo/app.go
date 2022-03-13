package vo

type Config struct {
	Mode           string     `json:"mode"`
	Debug          bool       `json:"debug"`
	Game           GameConfig `json:"game"`
	ManagerAddress []string   `json:"manager_address"`
	WhiteList      []string   `json:"white_list"`
}

type GameConfig struct {
	MinBalance          int   `json:"min_balance"`
	StartHour           int   `json:"start_hour"`
	EndHour             int   `json:"end_hour"`
	MinWithdraw         int64 `json:"min_withdraw"`
	VerifyWalletBalance bool  `json:"verify_wallet_balance"`
	OnlyWhiteList       bool  `json:"only_white_list"`
}
