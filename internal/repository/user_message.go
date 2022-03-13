package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserMessageRepo interface {
	Save(userMessage model.UserMessage) (model.UserMessage, error)
	PagingByAddressAndState(address string, state model.UserMessageState, page, size int) (int64, []model.UserMessage, error)
	FindById(id int64) (model.UserMessage, error)
}

func NewUserMessage(db *gorm.DB) UserMessageRepo {
	return &userMessageRepo{db}
}

type userMessageRepo struct {
	db *gorm.DB
}

func (r *userMessageRepo) Save(userMessage model.UserMessage) (model.UserMessage, error) {
	err := r.db.Save(&userMessage).Error
	if err != nil {
		logger.Errorf("save userMessage occur error")
		return model.UserMessage{}, err
	}
	return userMessage, nil
}

func (r *userMessageRepo) PagingByAddressAndState(address string, state model.UserMessageState, page, size int) (int64, []model.UserMessage, error) {

	Db := r.db.Model(model.UserMessage{})
	if address != "" {
		Db = Db.Where("address = ?", address)
	}
	if state != -1 {
		Db = Db.Where("state = ?", state)
	}

	var total int64
	err := Db.Count(&total).Error
	if err != nil {
		logger.WithError(err).Error("get userMessage total occur error")
		return 0, nil, err
	}

	userMessages := make([]model.UserMessage, 0)
	err = Db.Order("created_at desc").Offset(page * size).Limit(size).Find(&userMessages).Error
	if err != nil {
		logger.WithError(err).Error("get userMessage paging occur error")
		return 0, nil, err
	}
	return total, userMessages, nil
}

func (r *userMessageRepo) FindById(id int64) (model.UserMessage, error) {
	var userMessage model.UserMessage
	err := r.db.Model(model.UserMessage{}).Where("id = ?", id).First(&userMessage).Error
	if err != nil {
		logger.Errorf("query message by id occur error")
		return model.UserMessage{}, err
	}
	return userMessage, nil
}
