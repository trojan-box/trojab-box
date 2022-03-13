package model

type BonusWithdraw struct {
	Model
	Address string             `json:"address" gorm:"type:varchar(50)"`
	Bonus   int64              `json:"bonus"`
	Gas     int                `json:"gas"`
	State   BonusWithdrawState `json:"state"`
	Txhash  string             `json:"txhash" gorm:"type:varchar(120)"`
}

func (b *BonusWithdraw) TableName() string {
	return "t_bonus_withdraw"
}

type BonusWithdrawState int

const (
	BonusWithdrawStateSubmit     BonusWithdrawState = 0
	BonusWithdrawStateProcessing BonusWithdrawState = 1
	BonusWithdrawStateCompleted  BonusWithdrawState = 2
)
