package vo

type UserAuthReq struct {
	Address   string `json:"address"`
	Timestamp string `json:"timestamp"`
	Nonce     string `json:"nonce"`
	SignedMsg string `json:"signed_msg"`
}

type UserProfile struct {
	Address  string `json:"address"`
	NickName string `json:"nick_name"`
	Avatar   int    `json:"avatar"`
}

type ModifyUserProfileReq struct {
	NickName string `json:"nick_name"`
	Avatar   int    `json:"avatar"`
}
