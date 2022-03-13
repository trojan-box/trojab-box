package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserBonusRepo interface {
	Save(userBonus model.UserBonus) (model.UserBonus, error)
	Find(address string) (model.UserBonus, error)
	GetRankByTotalWinDesc(limit int) ([]model.UserBonus, error)
}

func NewUserBonusRepo(db *gorm.DB) UserBonusRepo {
	return &userBonusRepo{db}
}

type userBonusRepo struct {
	db *gorm.DB
}

func (r *userBonusRepo) Save(userBonus model.UserBonus) (model.UserBonus, error) {
	err := r.db.Save(&userBonus).Error
	if err != nil {
		logger.WithError(err).Errorf("save user bonus occur error")
		return model.UserBonus{}, err
	}
	return userBonus, nil
}

func (r *userBonusRepo) Find(address string) (model.UserBonus, error) {
	userBonus := model.UserBonus{}
	err := r.db.Where("address = ?", address).First(&userBonus).Error
	if err != nil {
		return model.UserBonus{}, err
	} else {
		return userBonus, nil
	}
}

func (r *userBonusRepo) GetRankByTotalWinDesc(limit int) ([]model.UserBonus, error) {
	userBonus := make([]model.UserBonus, 0)
	err := r.db.Order("total_win desc").Limit(limit).Find(&userBonus).Error
	if err != nil {
		logger.WithError(err).Errorf("find order by total_win occur error")
		return nil, err
	}
	return userBonus, nil
}
