package model

type DailyStats struct {
	Model
	Day               string  `json:"day" gorm:"uniqueIndex;type:varchar(20)"`
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

func (b *DailyStats) TableName() string {
	return "t_daily_stats"
}
