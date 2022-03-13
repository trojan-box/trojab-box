package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserYieldHourlyStatsRepo interface {
	Save(stats model.UserYieldHourlyStats) (model.UserYieldHourlyStats, error)
	Paging(page, size int) (int64, []model.UserYieldHourlyStats, error)
}

func NewUserYieldHourlyStats(db *gorm.DB) UserYieldHourlyStatsRepo {
	return &userYieldHourlyStatsRepo{db}
}

type userYieldHourlyStatsRepo struct {
	db *gorm.DB
}

func (r *userYieldHourlyStatsRepo) Save(stats model.UserYieldHourlyStats) (model.UserYieldHourlyStats, error) {
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "time"}},
		UpdateAll: true,
	}).Create(&stats).Error
	if err != nil {
		logger.WithError(err).Errorf("save user yield hourly stats occur error")
		return model.UserYieldHourlyStats{}, err
	}
	return stats, nil
}

func (r *userYieldHourlyStatsRepo) Paging(page, size int) (int64, []model.UserYieldHourlyStats, error) {
	var total int64
	err := r.db.Model(model.UserYieldHourlyStats{}).Count(&total).Error
	if err != nil {
		logger.WithError(err).Error("get user yield hourly stats total occur error")
		return 0, nil, err
	}
	stats := make([]model.UserYieldHourlyStats, 0)
	err = r.db.Model(model.UserYieldHourlyStats{}).Order("created_at desc").Offset(page * size).Limit(size).Find(&stats).Error
	if err != nil {
		logger.WithError(err).Error("get user yield hourly stats paging occur error")
		return 0, nil, err
	}
	return total, stats, nil
}
