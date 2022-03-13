package model

type UserBonus struct {
	Address   string `json:"address" gorm:"primarykey;type:varchar(50)"`
	Balance   int64  `json:"balance"`
	TotalWin  int64  `json:"total_win"`
	Freeze    int64  `json:"freeze"`
	CreatedAt int64
	UpdatedAt int64
}

func (u *UserBonus) TableName() string {
	return "t_user_bonus"
}
