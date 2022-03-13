package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LeaderboardRepo interface {
	DeleteByType(leaderboardType model.LeaderboardType) error
	Save(leaderboards []model.Leaderboard) ([]model.Leaderboard, error)
	GetLeaderboardByType(leaderboardType model.LeaderboardType) ([]model.Leaderboard, error)
}

func NewLeaderboard(db *gorm.DB) LeaderboardRepo {
	return &leaderboardRepo{db}
}

type leaderboardRepo struct {
	db *gorm.DB
}

func (r *leaderboardRepo) Save(leaderboards []model.Leaderboard) ([]model.Leaderboard, error) {
	err := r.db.Save(&leaderboards).Error
	if err != nil {
		logger.Errorf("save leaderboards occur error")
		return []model.Leaderboard{}, err
	}
	return leaderboards, nil
}

func (r *leaderboardRepo) GetLeaderboardByType(leaderboardType model.LeaderboardType) ([]model.Leaderboard, error) {
	leaderboards := make([]model.Leaderboard, 0)
	err := r.db.Model(model.Leaderboard{}).Where("type = ?", leaderboardType).Order("reward desc").Find(&leaderboards).Error
	if err != nil {
		logger.WithError(err).Error("get leaderboards by type occur error")
		return nil, err
	}
	return leaderboards, nil
}

func (r *leaderboardRepo) DeleteByType(leaderboardType model.LeaderboardType) error {
	err := r.db.Unscoped().Where("type = ?", leaderboardType).Delete(&model.Leaderboard{}).Error
	if err != nil {
		logger.WithError(err).Errorf("delete leaderboard by type occur error")
		return err
	}
	return nil
}
