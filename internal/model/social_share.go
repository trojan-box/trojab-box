package model

type SocialShare struct {
	Model
	Address        string             `json:"address" gorm:"type:varchar(50)"`
	Link           string             `json:"link"`
	Content        string             `json:"content"`
	ShareType      SocialShareType    `json:"share_type"`
	Bonus          int64              `json:"bonus"`
	State          BonusWithdrawState `json:"state"`
	Channel        SocialShareChannel `json:"channel"`
	Auditor        string             `json:"auditor"`
	AuditorAddress string             `json:"auditor_address" gorm:"type:varchar(50)"`
	Associate      int64              `json:"associate"`
	Accept         bool               `json:"accept"`
	Reply          string             `json:"reply"`
}

func (b *SocialShare) TableName() string {
	return "t_social_share"
}

type SocialShareState int

const (
	SocialShareStateSubmit     BonusWithdrawState = 0
	SocialShareStateProcessing BonusWithdrawState = 1
	SocialShareStateCompleted  BonusWithdrawState = 2
)

type SocialShareType int

const (
	SocialShareTypeWithdraw SocialShareType = 1
	SocialShareTypeCommon   SocialShareType = 2
	SocialShareTypeReport   SocialShareType = 3
)

type SocialShareChannel int

const (
	SocialShareChannelGate     SocialShareChannel = 1
	SocialShareChannelWeibo    SocialShareChannel = 2
	SocialShareChannelTwitter  SocialShareChannel = 3
	SocialShareChannelReddit   SocialShareChannel = 4
	SocialShareChannelFacebook SocialShareChannel = 5
	SocialShareChannelWebsite  SocialShareChannel = 6
)

var SocialShareChannelName = map[SocialShareChannel]string{
	SocialShareChannelGate:     "Gate",
	SocialShareChannelWeibo:    "Weibo",
	SocialShareChannelTwitter:  "Twitter",
	SocialShareChannelReddit:   "Reddit",
	SocialShareChannelFacebook: "Facebook",
	SocialShareChannelWebsite:  "Website",
}
