package vo

import "github.com/aresprotocols/trojan-box/internal/model"

type SocialShare struct {
	ID             int64                    `json:"id"`
	NickName       string                   `json:"nick_name"`
	Address        string                   `json:"address" gorm:"type:varchar(50)"`
	Link           string                   `json:"link"`
	Content        string                   `json:"content"`
	ShareType      model.SocialShareType    `json:"share_type"`
	Channel        model.SocialShareChannel `json:"channel"`
	Bonus          int64                    `json:"bonus"`
	State          model.BonusWithdrawState `json:"state"`
	Auditor        string                   `json:"auditor"`
	AuditorAddress string                   `json:"auditor_address" gorm:"type:varchar(50)"`
	ApplyTime      int                      `json:"apply_time" copier:"CreatedAt"`
	Reply          string                   `json:"reply"`
}

type AddSocialShareReq struct {
	Link      string                   `json:"link"`
	Content   string                   `json:"content"`
	ShareType model.SocialShareType    `json:"share_type"`
	Channel   model.SocialShareChannel `json:"channel"`
}

type ProcessSocialShareReq struct {
	ID      int64  `json:"id"`
	Bonus   int64  `json:"bonus"`
	Auditor string `json:"auditor"`
	Accept  bool   `json:"accept"`
	Reply   string `json:"reply"`
}
