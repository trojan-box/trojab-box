package repository

import (
	"errors"
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/vo"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GameSessionRepo interface {
	FindByAddressAndSession(address string, session string) (model.GameSession, error)
	FindByAddressAndId(address string, id int64) (model.GameSession, error)
	Save(gameSession model.GameSession) (model.GameSession, error)
	PagingByAddress(address string, page, size int) (int64, []model.GameSession, error)
	GetRankByReward(startTime, endTime int64, limit int) ([]vo.NewStarRecord, error)
	GetWinBigBonusByTime(startTime, endTime int64) ([]model.GameSession, error)
	GetPartInAddressCount(startTime, endTime int64) (int64, error)
	GetPartInCount(startTime, endTime int64) (int64, error)
	GetRewardAmount(startTime, endTime int64) (int64, error)
	GetMaxMinSingleRewardSum(startTime, endTime int64, isMax bool) (int64, error)
	GetTotalRewardAmount() (int64, error)
	GetTotalCount() (int64, error)
	LatestWinBigBonus(address string) (model.GameSession, error)
	CountAfterTime(address string, timestamp int64) (int64, error)
}

func NewGameSessionRepo(db *gorm.DB) GameSessionRepo {
	return &gameSessionRepo{db}
}

type gameSessionRepo struct {
	db *gorm.DB
}

func (r *gameSessionRepo) Save(gameSession model.GameSession) (model.GameSession, error) {
	err := r.db.Save(&gameSession).Error
	if err != nil {
		logger.Errorf("save game session occur error")
		return model.GameSession{}, err
	}
	return gameSession, nil
}

func (r *gameSessionRepo) FindByAddressAndSession(address string, session string) (model.GameSession, error) {
	gameSession := model.GameSession{}
	err := r.db.Model(gameSession).Where("address = ? and session = ?", address, session).First(&gameSession).Error
	if err != nil {
		logger.WithField("address", address).WithField("session", session).WithError(err).Errorf("find game session by address and session occur error")
		return model.GameSession{}, err
	}
	return gameSession, nil
}

func (r *gameSessionRepo) PagingByAddress(address string, page, size int) (int64, []model.GameSession, error) {

	log := logger.WithField("address", address)

	Db := r.db.Model(model.GameSession{})
	if address != "" {
		Db = Db.Where("address = ?", address)
	}

	var total int64
	err := Db.Count(&total).Error
	if err != nil {
		log.WithError(err).Error("get game session total occur error")
		return 0, nil, err
	}

	games := make([]model.GameSession, 0)
	err = Db.Order("created_at desc").Offset(page * size).Limit(size).Find(&games).Error
	if err != nil {
		log.WithError(err).Error("get game session paging occur error")
		return 0, nil, err
	}
	return total, games, nil
}

func (r *gameSessionRepo) FindByAddressAndId(address string, id int64) (model.GameSession, error) {
	gameSession := model.GameSession{}
	err := r.db.Model(gameSession).Where("address = ? and id = ?", address, id).First(&gameSession).Error
	if err != nil {
		logger.WithField("address", address).WithField("session", id).Errorf("find game session by address and id occur error")
		return model.GameSession{}, err
	}
	return gameSession, nil
}

func (r *gameSessionRepo) GetRankByReward(startTime, endTime int64, limit int) ([]vo.NewStarRecord, error) {
	newStarRecords := make([]vo.NewStarRecord, 0)
	err := r.db.Model(model.GameSession{}).Select("address,IFNULL(sum(bonus), 0) as reward").
		Where("created_at>= ? and created_at <= ?", startTime, endTime).
		Group("address").Order("reward desc").Limit(limit).Find(&newStarRecords).Error
	if err != nil {
		logger.WithError(err).Errorf("get new start from game session occur error")
		return nil, err
	}
	return newStarRecords, nil
}

func (r *gameSessionRepo) GetWinBigBonusByTime(startTime, endTime int64) ([]model.GameSession, error) {
	games := make([]model.GameSession, 0)
	err := r.db.Model(model.GameSession{}).
		Where("created_at>= ? and created_at <= ? and bonus_level < 9", startTime, endTime).
		Order("bonus_level desc").Find(&games).Error
	if err != nil {
		logger.WithError(err).Error("get win big bonus occur error")
		return nil, err
	}
	return games, nil
}

func (r *gameSessionRepo) GetPartInAddressCount(startTime, endTime int64) (int64, error) {
	var count int64

	Db := r.db.Model(model.GameSession{})
	if startTime != 0 {
		Db = Db.Where("created_at>= ?", startTime)
	}
	if endTime != 0 {
		Db = Db.Where("created_at <= ?", endTime)
	}
	err := Db.Distinct("address").Count(&count).Error
	if err != nil {
		logger.WithError(err).Error("get part in address count occur error")
		return 0, err
	}
	return count, nil
}
func (r *gameSessionRepo) GetPartInCount(startTime, endTime int64) (int64, error) {
	var count int64
	err := r.db.Model(model.GameSession{}).
		Where("created_at>= ? and created_at <= ?", startTime, endTime).Count(&count).Error
	if err != nil {
		logger.WithError(err).Error("get part in count occur error")
		return 0, err
	}
	return count, nil
}
func (r *gameSessionRepo) GetRewardAmount(startTime, endTime int64) (int64, error) {
	var sum int64
	err := r.db.Model(model.GameSession{}).Select("IFNULL(sum(bonus), 0)").
		Where("created_at>= ? and created_at <= ?", startTime, endTime).Find(&sum).Error
	if err != nil {
		logger.WithError(err).Error("get sum bonus occur error")
		return 0, err
	}
	return sum, nil
}

func (r *gameSessionRepo) GetMaxMinSingleRewardSum(startTime, endTime int64, isMax bool) (int64, error) {
	var sum int64
	var order = "reward asc"
	if isMax {
		order = "reward desc"
	}
	err := r.db.Model(model.GameSession{}).Select("IFNULL(sum(bonus), 0) as reward").
		Where("created_at>= ? and created_at <= ?", startTime, endTime).
		Group("address").Order(order).Limit(1).Find(&sum).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			logger.WithError(err).Errorf("get new start from game session occur error")
			return 0, err
		}
	}
	return sum, nil
}

func (r *gameSessionRepo) GetTotalRewardAmount() (int64, error) {
	var sum int64
	err := r.db.Model(model.GameSession{}).Select("IFNULL(sum(bonus), 0)").Find(&sum).Error
	if err != nil {
		logger.WithError(err).Error("get total sum bonus occur error")
		return 0, err
	}
	return sum, nil
}
func (r *gameSessionRepo) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(model.GameSession{}).Count(&count).Error
	if err != nil {
		logger.WithError(err).Error("get total count occur error")
		return 0, err
	}
	return count, nil
}

func (r *gameSessionRepo) LatestWinBigBonus(address string) (model.GameSession, error) {
	gameSession := model.GameSession{}
	err := r.db.Model(gameSession).Where("address = ? and bonus_level < 9", address).Order("created_at desc").First(&gameSession).Error
	if err != nil {
		logger.WithField("address", address).WithError(err).Errorf("find latest win big bonus occur error")
		return model.GameSession{}, err
	}
	return gameSession, nil
}

func (r *gameSessionRepo) CountAfterTime(address string, timestamp int64) (int64, error) {
	var count int64
	err := r.db.Model(model.GameSession{}).Where("address = ? and created_at > ?", address, timestamp).Count(&count).Error
	if err != nil {
		logger.WithField("address", address).WithError(err).Errorf("count after time occur err")
		return 0, err
	}
	return count, nil
}
