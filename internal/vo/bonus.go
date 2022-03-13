package vo

import (
	"github.com/aresprotocols/trojan-box/internal/model"
)

type UserBonus struct {
	Address  string `json:"address"`
	Balance  int64  `json:"balance"`
	TodayWin int64  `json:"today_win"`
	TotalWin int64  `json:"total_win"`
	Freeze   int64  `json:"freeze"`
}

type BonusHistory struct {
	Address    string                 `json:"address"`
	Bonus      int64                  `json:"bonus"`
	Type       model.BonusRecordType  `json:"type"`
	Associate  int64                  `json:"associate"`
	State      model.BonusRecordState `json:"state"`
	CreateTime int64                  `json:"create_time" copier:"CreatedAt"`
}

type WithdrawBonusApplyReq struct {
	Address   string `json:"address"`
	Nonce     string `json:"nonce"`
	Timestamp string `json:"timestamp"`
	Bonus     int64  `json:"bonus"`
	SignedMsg string `json:"signed_msg"`
}
