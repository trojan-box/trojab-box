package vo

type DailyStats struct {
	Day               string  `json:"day"`
	NewAddress        int64   `json:"new_address"`
	PartInAddress     int64   `json:"part_in_address"`
	PartInCount       int64   `json:"part_in_count"`
	RewardAmount      int64   `json:"reward_amount"`
	WithdrawAmount    int64   `json:"withdraw_amount"`
	OpenedBigReward   int     `json:"opened_big_reward"`
	UnopenedBigReward int     `json:"unopened_big_reward"`
	SingleRewardMax   int64   `json:"single_reward_max"`
	SingleRewardMin   int64   `json:"single_reward_min"`
	StakingAmount     int64   `json:"staking_amount"`
	GameRewardAmount  int64   `json:"game_reward_amount"`
	ShareRewardAmount int64   `json:"share_reward_amount"`
	AnnualYieldRate   float64 `json:"annual_yield_rate"`
}

type TotalStats struct {
	SignedAddress     int64   `json:"signed_address"`
	PartInAddress     int64   `json:"part_in_address"`
	PartInCount       int64   `json:"part_in_count"`
	RewardAmount      int64   `json:"reward_amount"`
	WithdrawAmount    int64   `json:"withdraw_amount"`
	GameRewardAmount  int64   `json:"game_reward_amount"`
	ShareRewardAmount int64   `json:"share_reward_amount"`
	StakingAmount     int64   `json:"staking_amount"`
	AvgAPR            float64 `json:"avg_apr"`
}

type UserYieldHourlyStats struct {
	Time         string  `json:"time"`
	Balance      int64   `json:"balance"`
	Reward       int64   `json:"reward"`
	AnnuallyRate float64 `json:"annually_rate"`
}
