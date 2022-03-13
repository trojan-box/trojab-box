package model

type UserYieldHourlyStats struct {
	Model
	Time         string  `json:"time" gorm:"uniqueIndex;type:varchar(20)"`
	Balance      int64   `json:"balance"`
	Reward       int64   `json:"reward"`
	AnnuallyRate float64 `json:"annually_rate"`
}

func (b *UserYieldHourlyStats) TableName() string {
	return "t_user_yield_hourly_stats"
}
