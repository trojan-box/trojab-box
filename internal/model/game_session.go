package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
)

type GameSession struct {
	Model
	Address     string              `json:"address" grom:"index;type:varchar(50)"`
	Session     string              `json:"session" gorm:"type:varchar(40)"`
	ChosenIndex int                 `json:"chosen_index"`
	Bonus       int64               `json:"bonus"`
	BonusLevel  constant.BonusLevel `json:"bonus_level"`
	Cards       GodSlice            `json:"cards" gorm:"type:text"`
	CardsBonus  IntSlice            `json:"cards_bonus" gorm:"type:text"`
}

func (g *GameSession) TableName() string {
	return "t_game_session"
}

type IntSlice []int64

func (j *IntSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal str value:", value))
	}
	var result []int64
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

func (j IntSlice) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

type GodSlice []constant.God

func (j *GodSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal str value:", value))
	}
	var result []constant.God
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

func (j GodSlice) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}
