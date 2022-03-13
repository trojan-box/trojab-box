package vo

type LeaderboardResp struct {
	NickName string `json:"nick_name"`
	Avatar   int    `json:"avatar"`
	Address  string `json:"address"`
	Reward   int    `json:"reward"`
}
