package vo

import "github.com/aresprotocols/trojan-box/internal/model"

type UserMessage struct {
	ID      int64                  `json:"id"`
	Content string                 `json:"content"`
	State   model.UserMessageState `json:"state"`
}

type ReadMessage struct {
	ID int64 `json:"id"`
}
