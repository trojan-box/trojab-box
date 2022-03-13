package model

type Leaderboard struct {
	Model
	Address string          `json:"address" gorm:"type:varchar(50)"`
	Reward  int64           `json:"reward"`
	Type    LeaderboardType `json:"type"`
}

func (l *Leaderboard) TableName() string {
	return "t_leaderboard"
}

type LeaderboardType int

const (
	LeaderboardTypeNewStar        LeaderboardType = 1
	LeaderboardTypeSeasonChampion LeaderboardType = 2
)
