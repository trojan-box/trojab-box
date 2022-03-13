package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BonusRecordRepo interface {
	Save(record model.BonusRecord) (model.BonusRecord, error)
	PagingByAddress(address string, page, size int, bonusRecordType int) (int64, []model.BonusRecord, error)
	GetByTimeAndAddress(startTime, endTime int64, address string) ([]model.BonusRecord, error)
	FindByTypeAndAssociate(recordType model.BonusRecordType, associate int64) (model.BonusRecord, error)
}

func NewBonusRecordRepo(db *gorm.DB) BonusRecordRepo {
	return &bonusRecordRepo{db}
}

type bonusRecordRepo struct {
	db *gorm.DB
}

func (r *bonusRecordRepo) Save(record model.BonusRecord) (model.BonusRecord, error) {
	err := r.db.Save(&record).Error
	if err != nil {
		logger.Errorf("save bonus record occur error")
		return model.BonusRecord{}, err
	}
	return record, nil
}

func (r *bonusRecordRepo) PagingByAddress(address string, page, size int, bonusRecordType int) (int64, []model.BonusRecord, error) {

	log := logger.WithField("address", address)

	Db := r.db.Model(model.BonusRecord{}).Where("address = ?", address)
	if bonusRecordType != -1 {
		Db = Db.Where("type = ?", bonusRecordType)
	}
	var total int64
	err := Db.Count(&total).Error
	if err != nil {
		log.WithError(err).Error("get bonus record total occur error")
		return 0, nil, err
	}
	records := make([]model.BonusRecord, 0)
	err = Db.Order("created_at desc").Offset(page * size).Limit(size).Find(&records).Error
	if err != nil {
		log.WithError(err).Error("get bonus record paging occur error")
		return 0, nil, err
	}
	return total, records, nil
}

func (r *bonusRecordRepo) GetByTimeAndAddress(startTime, endTime int64, address string) ([]model.BonusRecord, error) {
	games := make([]model.BonusRecord, 0)
	err := r.db.Model(model.BonusRecord{}).
		Where("created_at>= ? and created_at <= ? and address = ?", startTime, endTime, address).Find(&games).Error
	if err != nil {
		logger.WithError(err).Error("get by time and address occur error")
		return nil, err
	}
	return games, nil
}
func (r *bonusRecordRepo) FindByTypeAndAssociate(recordType model.BonusRecordType, associate int64) (model.BonusRecord, error) {
	record := model.BonusRecord{}
	err := r.db.Model(model.BonusRecord{}).
		Where("type = ? and associate = ?", recordType, associate).First(&record).Error
	if err != nil {
		logger.WithError(err).Error("get record by type and associate occur error")
		return model.BonusRecord{}, err
	}
	return record, nil
}
