package model

type UserWalletBalance struct {
	Address   string `json:"address" gorm:"primarykey;type:varchar(50)"`
	Time      string `json:"time" gorm:"primarykey;type:varchar(20)"`
	Balance   int64  `json:"balance"`
	CreatedAt int64
	UpdatedAt int64
}

func (u *UserWalletBalance) TableName() string {
	return "t_user_wallet_balance"
}
