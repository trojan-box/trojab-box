package model

type Nonce struct {
	Model
	Address string `json:"address" gorm:"type:varchar(50)"`
	Nonce   string `json:"nonce" gorm:"type:varchar(10)"`
}

func (b *Nonce) TableName() string {
	return "t_nonce"
}
