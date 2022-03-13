package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BroadcastRepo interface {
	Save(broadcast model.Broadcast) (model.Broadcast, error)
	PagingByAddress(page, size int) (int64, []model.Broadcast, error)
	GetLatest() (model.Broadcast, error)
}

func NewBroadcast(db *gorm.DB) BroadcastRepo {
	return &broadcastRepo{db}
}

type broadcastRepo struct {
	db *gorm.DB
}

func (r *broadcastRepo) Save(broadcast model.Broadcast) (model.Broadcast, error) {
	err := r.db.Save(&broadcast).Error
	if err != nil {
		logger.Errorf("save broadcast occur error")
		return model.Broadcast{}, err
	}
	return broadcast, nil
}

func (r *broadcastRepo) GetLatest() (model.Broadcast, error) {
	broadcast := model.Broadcast{}
	err := r.db.Order("created_at desc").First(&broadcast).Error
	if err != nil {
		logger.WithError(err).Errorf("query latest broadcast occur err")
		return model.Broadcast{}, err
	}
	return broadcast, err
}

func (r *broadcastRepo) PagingByAddress(page, size int) (int64, []model.Broadcast, error) {

	var total int64
	err := r.db.Model(model.Broadcast{}).Count(&total).Error
	if err != nil {
		logger.WithError(err).Error("get broadcast total occur error")
		return 0, nil, err
	}

	broadcasts := make([]model.Broadcast, 0)
	err = r.db.Model(model.Broadcast{}).Order("created_at desc").Offset(page * size).Limit(size).Find(&broadcasts).Error
	if err != nil {
		logger.WithError(err).Error("get broadcast paging occur error")
		return 0, nil, err
	}
	return total, broadcasts, nil
}
