package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BonusWithdrawRepo interface {
	Save(bonusWithdraw model.BonusWithdraw) (model.BonusWithdraw, error)
	GetSumByTime(startTime, endTime int64) (int64, error)
	GetTotalWithdrawAmount() (int64, error)
	PagingByStateAndAddress(state model.BonusWithdrawState, address string, page, size int) (int64, []model.BonusWithdraw, error)
	FindById(id int64) (model.BonusWithdraw, error)
	GetTotalCount() (int64, error)
	FindLastByAddress(address string) (model.BonusWithdraw, error)
}

func NewBonusWithdrawRepo(db *gorm.DB) BonusWithdrawRepo {
	return &bonusWithdrawRepo{db}
}

type bonusWithdrawRepo struct {
	db *gorm.DB
}

func (r *bonusWithdrawRepo) Save(bonusWithdraw model.BonusWithdraw) (model.BonusWithdraw, error) {
	err := r.db.Save(&bonusWithdraw).Error
	if err != nil {
		logger.Errorf("save bonus withdraw occur error")
		return model.BonusWithdraw{}, err
	}
	return bonusWithdraw, nil
}

func (r *bonusWithdrawRepo) FindById(id int64) (model.BonusWithdraw, error) {
	var bonusWithdraw model.BonusWithdraw
	err := r.db.Model(model.BonusWithdraw{}).Where("id = ?", id).First(&bonusWithdraw).Error
	if err != nil {
		logger.Errorf("query bonus withdraw by id occur error")
		return model.BonusWithdraw{}, err
	}
	return bonusWithdraw, nil
}

func (r *bonusWithdrawRepo) GetSumByTime(startTime, endTime int64) (int64, error) {
	var sum int64
	err := r.db.Model(model.BonusWithdraw{}).Select("IFNULL(sum(bonus), 0) as reward").
		Where("created_at>= ? and created_at <= ?", startTime, endTime).First(&sum).Error
	if err != nil {
		logger.WithError(err).Errorf("get sum from bonus withdraw occur error")
		return 0, err
	}
	return sum, nil
}
func (r *bonusWithdrawRepo) GetTotalWithdrawAmount() (int64, error) {
	var sum int64
	err := r.db.Model(model.BonusWithdraw{}).Select("IFNULL(sum(bonus), 0) as reward").First(&sum).Error
	if err != nil {
		logger.WithError(err).Errorf("get total sum from bonus withdraw occur error")
		return 0, err
	}
	return sum, nil
}

func (r *bonusWithdrawRepo) PagingByStateAndAddress(state model.BonusWithdrawState, address string, page, size int) (int64, []model.BonusWithdraw, error) {
	log := logger.WithField("state", state)
	Db := r.db.Model(model.BonusWithdraw{})
	if state != -1 {
		Db = Db.Where("state = ?", state)
	}
	if address != "" {
		Db = Db.Where("address = ?", address)
	}
	var total int64
	err := Db.Count(&total).Error
	if err != nil {
		log.WithError(err).Error("get bonus withdraw by state total occur error")
		return 0, nil, err
	}
	bonusWithdraws := make([]model.BonusWithdraw, 0)
	err = Db.Order("created_at desc").Offset(page * size).Limit(size).Find(&bonusWithdraws).Error
	if err != nil {
		log.WithError(err).Error("get bonus withdraw by state paging occur error")
		return 0, nil, err
	}
	return total, bonusWithdraws, nil
}

func (r *bonusWithdrawRepo) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(model.BonusWithdraw{}).Count(&count).Error
	if err != nil {
		logger.WithError(err).Error("get total count occur error")
		return 0, err
	}
	return count, nil
}

func (r *bonusWithdrawRepo) FindLastByAddress(address string) (model.BonusWithdraw, error) {
	var bonusWithdraw model.BonusWithdraw
	err := r.db.Model(model.BonusWithdraw{}).Where("address = ?", address).Order("created_at desc").First(&bonusWithdraw).Error
	if err != nil {
		logger.Errorf("query last  bonus withdraw by address occur error")
		return model.BonusWithdraw{}, err
	}
	return bonusWithdraw, nil
}
