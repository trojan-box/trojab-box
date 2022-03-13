package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DailyStatsRepo interface {
	Save(dailyStats model.DailyStats) (model.DailyStats, error)
	GetByDay(day string) (model.DailyStats, error)
	Paging(page, size int) (int64, []model.DailyStats, error)
	GetSumAnnualYieldRateAndCount() (int64, float64, error)
}

func NewDailyStats(db *gorm.DB) DailyStatsRepo {
	return &dailyStatsRepo{db}
}

type dailyStatsRepo struct {
	db *gorm.DB
}

func (r *dailyStatsRepo) Save(dailyStats model.DailyStats) (model.DailyStats, error) {
	err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "day"}},
		UpdateAll: true,
	}).Create(&dailyStats).Error
	if err != nil {
		logger.WithError(err).Errorf("save dailyStats occur error")
		return model.DailyStats{}, err
	}
	return dailyStats, nil
}

func (r *dailyStatsRepo) GetByDay(day string) (model.DailyStats, error) {
	var dailyStats model.DailyStats
	err := r.db.Model(model.DailyStats{}).Where("day = ?", day).First(&dailyStats).Error
	if err != nil {
		logger.WithError(err).Error("get daily stats by type occur error")
		return model.DailyStats{}, err
	}
	return dailyStats, nil
}

func (r *dailyStatsRepo) Paging(page, size int) (int64, []model.DailyStats, error) {
	var total int64
	err := r.db.Model(model.DailyStats{}).Count(&total).Error
	if err != nil {
		logger.WithError(err).Error("get daily stats total occur error")
		return 0, nil, err
	}
	stats := make([]model.DailyStats, 0)
	err = r.db.Model(model.DailyStats{}).Order("created_at desc").Offset(page * size).Limit(size).Find(&stats).Error
	if err != nil {
		logger.WithError(err).Error("get daily stats paging occur error")
		return 0, nil, err
	}
	return total, stats, nil
}

func (r *dailyStatsRepo) GetSumAnnualYieldRateAndCount() (int64, float64, error) {
	Db := r.db.Model(model.DailyStats{})
	Db = Db.Where("staking_amount > ?", 0)
	var total int64
	err := Db.Count(&total).Error
	if err != nil {
		logger.WithError(err).Error("get yield count occur error")
		return 0, 0, err
	}
	var sum float64
	err = Db.Select("IFNULL(sum(annual_yield_rate), 0.0) as rate").Find(&sum).Error
	if err != nil {
		logger.WithError(err).Error("get sum yield occur error")
		return 0, 0, err
	}
	return total, sum, nil
}
