package vo

import "github.com/aresprotocols/trojan-box/internal/pkg/constant"

type PlayGameReq struct {
	Address     string         `json:"address"`
	Nonce       string         `json:"nonce"`
	Timestamp   string         `json:"timestamp"`
	Cards       []constant.God `json:"cards"`
	ChosenIndex int            `json:"chosen"`
	SignedMsg   string         `json:"signed_msg"`
}

type PlayGameResp struct {
	ID      int64  `json:"id"`
	Session string `json:"session"`
	Bonus   int64  `json:"bonus"`
}

type GameSession struct {
	ID          int64          `json:"id"`
	Address     string         `json:"address" `
	Session     string         `json:"session" `
	ChosenIndex int            `json:"chosen_index"`
	Bonus       int64          `json:"bonus"`
	PlayTime    int64          `json:"play_time" copier:"CreatedAt"`
	Cards       []constant.God `json:"cards"`
	CardsBonus  []int64        `json:"cards_bonus"`
}

type GameHistory struct {
	ID          int64          `json:"id"`
	Address     string         `json:"address" `
	Session     string         `json:"session" `
	Bonus       int64          `json:"bonus"`
	PlayTime    int            `json:"play_time" copier:"CreatedAt"`
	ChosenIndex int            `json:"chosen_index"`
	Cards       []constant.God `json:"cards"`
	CardsBonus  []int64        `json:"cards_bonus"`
}
type GameHistories struct {
	ID          int64          `json:"id"`
	NickName    string         `json:"nick_name"`
	Address     string         `json:"address" `
	Session     string         `json:"session" `
	Bonus       int64          `json:"bonus"`
	PlayTime    int            `json:"play_time" copier:"CreatedAt"`
	ChosenIndex int            `json:"chosen_index"`
	Cards       []constant.God `json:"cards"`
	CardsBonus  []int64        `json:"cards_bonus"`
}

type NewStarRecord struct {
	Address string `json:"address" `
	Reward  int64  `json:"reward"`
}
