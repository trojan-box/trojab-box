package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SocialShareRepo interface {
	Save(socialShare model.SocialShare) (model.SocialShare, error)
	PagingAndFilter(state model.SocialShareState, address string,
		shareType model.SocialShareType, channel model.SocialShareChannel, auditorAddress string, beginTime, endTime int64, page, size int) (int64, []model.SocialShare, error)
	FindById(id int64) (model.SocialShare, error)
	GetRewardAmount(startTime, endTime int64) (int64, error)
	GetTotalRewardAmount() (int64, error)
	FindByAddressAndTypeAndAssociate(address string, shareType model.SocialShareType, associate int64) (model.SocialShare, error)
	FindByAddressAndTypeAndChannel(address string, shareType model.SocialShareType, channel model.SocialShareChannel,
		beginTime, endTime int64) (model.SocialShare, error)
}

func NewSocialShareRepo(db *gorm.DB) SocialShareRepo {
	return &socialShareRepo{db}
}

type socialShareRepo struct {
	db *gorm.DB
}

func (r *socialShareRepo) Save(socialShare model.SocialShare) (model.SocialShare, error) {
	err := r.db.Save(&socialShare).Error
	if err != nil {
		logger.Errorf("save social share occur error")
		return model.SocialShare{}, err
	}
	return socialShare, nil
}

func (r *socialShareRepo) FindById(id int64) (model.SocialShare, error) {
	var socialShare model.SocialShare
	err := r.db.Model(model.SocialShare{}).Where("id = ?", id).First(&socialShare).Error
	if err != nil {
		logger.Errorf("query social share by id occur error")
		return model.SocialShare{}, err
	}
	return socialShare, nil
}

func (r *socialShareRepo) PagingAndFilter(state model.SocialShareState, address string,
	shareType model.SocialShareType, channel model.SocialShareChannel, auditorAddress string, beginTime, endTime int64, page, size int) (int64, []model.SocialShare, error) {

	log := logger.WithField("state", state)

	Db := r.db.Model(model.SocialShare{})

	if state != -1 {
		Db = Db.Where("state = ?", state)
	}
	if shareType != -1 {
		Db = Db.Where("share_type = ?", shareType)
	}
	if channel != -1 {
		Db = Db.Where("channel = ?", channel)
	}
	if address != "" {
		Db = Db.Where("address = ?", address)
	}
	if beginTime != -1 {
		Db = Db.Where("created_at >= ?", beginTime)
	}
	if endTime != -1 {
		Db = Db.Where("created_at <= ?", endTime)
	}
	if auditorAddress != "" {
		Db = Db.Where("auditor_address = ?", auditorAddress)
	}

	var total int64
	err := Db.Count(&total).Error
	if err != nil {
		log.WithError(err).Error("get social share by state total occur error")
		return 0, nil, err
	}
	socialShares := make([]model.SocialShare, 0)
	err = Db.Order("created_at desc").Offset(page * size).Limit(size).Find(&socialShares).Error
	if err != nil {
		log.WithError(err).Error("get social share by state and address paging occur error")
		return 0, nil, err
	}
	return total, socialShares, nil
}

func (r *socialShareRepo) GetRewardAmount(startTime, endTime int64) (int64, error) {
	var sum int64
	err := r.db.Model(model.SocialShare{}).Select("IFNULL(sum(bonus), 0)").
		Where("created_at>= ? and created_at <= ?", startTime, endTime).Find(&sum).Error
	if err != nil {
		logger.WithError(err).Error("get sum bonus occur error")
		return 0, err
	}
	return sum, nil
}
func (r *socialShareRepo) GetTotalRewardAmount() (int64, error) {
	var sum int64
	err := r.db.Model(model.SocialShare{}).Select("IFNULL(sum(bonus), 0)").Find(&sum).Error
	if err != nil {
		logger.WithError(err).Error("get total sum bonus occur error")
		return 0, err
	}
	return sum, nil
}

func (r *socialShareRepo) FindByAddressAndTypeAndAssociate(address string, shareType model.SocialShareType,
	associate int64) (model.SocialShare, error) {

	var socialShare model.SocialShare
	err := r.db.Model(model.SocialShare{}).Where("address = ? and share_type = ? and associate = ?",
		address, shareType, associate).First(&socialShare).Error
	if err != nil {
		logger.Errorf("query social share by address,shareType,associate occur error")
		return model.SocialShare{}, err
	}
	return socialShare, nil
}

func (r *socialShareRepo) FindByAddressAndTypeAndChannel(address string, shareType model.SocialShareType,
	channel model.SocialShareChannel, beginTime, endTime int64) (model.SocialShare, error) {

	var socialShare model.SocialShare
	err := r.db.Model(model.SocialShare{}).Where("address = ? and share_type = ? and channel = ? and created_at >= ? and created_at <= ?",
		address, shareType, channel, beginTime, endTime).First(&socialShare).Error
	if err != nil {
		logger.Errorf("query social share by address,shareType,channel occur error")
		return model.SocialShare{}, err
	}
	return socialShare, nil
}
