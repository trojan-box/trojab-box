package model

type User struct {
	Address   string `json:"address" gorm:"primarykey;type:varchar(50)"`
	NickName  string `json:"nick_name" `
	Avatar    int    `json:"avatar"`
	CreatedAt int64
	UpdatedAt int64
}

func (u *User) TableName() string {
	return "t_user"
}
