package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Broadcast struct {
	Model
	Address      string          `json:"address" gorm:"type:varchar(50)"`
	TemplateKey  string          `json:"template_key"`
	TemplateData TemplateDataMap `json:"template_data" gorm:"type:text"`
}

func (b *Broadcast) TableName() string {
	return "t_broadcast"
}

type TemplateDataMap map[string]interface{}

func (j *TemplateDataMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal str value:", value))
	}
	var result map[string]interface{}
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

func (j TemplateDataMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}
