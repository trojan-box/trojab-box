package vo

import "github.com/aresprotocols/trojan-box/internal/model"

type BonusWithdraw struct {
	ID        int64                    `json:"id"`
	Address   string                   `json:"address"`
	NickName  string                   `json:"nick_name"`
	Bonus     int64                    `json:"bonus"`
	Gas       int                      `json:"gas"`
	State     model.BonusWithdrawState `json:"state"`
	Txhash    string                   `json:"txhash" `
	ApplyTime int                      `json:"apply_time" copier:"CreatedAt"`
}

type ProcessWithdrawBonusReq struct {
	ID     int64  `json:"id"`
	Txhash string `json:"txhash" `
}
