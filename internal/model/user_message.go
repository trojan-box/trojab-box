package model

type UserMessage struct {
	Model
	Address      string           `json:"address" gorm:"type:varchar(50)"`
	TemplateKey  string           `json:"template_key"`
	TemplateData TemplateDataMap  `json:"template_data" gorm:"type:text"`
	IsTemplate   bool             `json:"is_template"`
	State        UserMessageState `json:"state"`
}

func (b *UserMessage) TableName() string {
	return "t_user_message"
}

type UserMessageState int

const (
	UserMessageStateUnread UserMessageState = 1
	UserMessageStateRead   UserMessageState = 2
)
