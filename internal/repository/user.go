package repository

import (
	"errors"
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepo interface {
	InsertUser(user model.User) error
	GetUserByAddress(address string) (model.User, error)
	UpdateUserProfile(user model.User) error
	GetNewAddressCount(startTime, endTime int64) (int64, error)
	GetUserTotalCount() (int64, error)
	GetUserByAddresses(addresses []string) ([]model.User, error)
}

func NewUser(db *gorm.DB) UserRepo {
	return &userRepo{db}
}

type userRepo struct {
	db *gorm.DB
}

func (u *userRepo) InsertUser(user model.User) error {
	log := logger.WithField("address", user.Address)
	err := u.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}},
		DoNothing: true,
	}).Create(&user).Error
	if err != nil {
		log.WithError(err).Errorf("insert user occur error")
		return err
	}
	return nil
}

func (u *userRepo) GetUserByAddress(address string) (model.User, error) {
	log := logger.WithField("address", address)
	user := model.User{}
	err := u.db.Where("address = ?", address).First(&user).Error
	if err != nil {
		log.WithError(err).Errorf("query user by address occur err")
		return model.User{}, err
	}
	return user, err
}

func (u *userRepo) UpdateUserProfile(user model.User) error {
	log := logger.WithField("address", user.Address)
	err := u.db.Model(model.User{}).Where("address = ?", user.Address).Updates(model.User{
		NickName: user.NickName,
		Avatar:   user.Avatar,
	}).Error
	if err != nil {
		log.WithError(err).Errorf("update user profile occur error")
		return err
	}
	return nil
}

func (u *userRepo) GetUserTotalCount() (int64, error) {
	var count int64
	err := u.db.Model(model.User{}).Distinct("address").Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			logger.WithError(err).Error("get total count occur error")
			return 0, err
		}
	}
	return count, nil
}

func (u *userRepo) GetNewAddressCount(startTime, endTime int64) (int64, error) {
	var count int64
	err := u.db.Model(model.User{}).
		Where("created_at>= ? and created_at <= ?", startTime, endTime).Distinct("address").Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		} else {
			logger.WithError(err).Error("get new address occur error")
			return 0, err
		}
	}
	return count, nil
}

func (u *userRepo) GetUserByAddresses(addresses []string) ([]model.User, error) {
	users := make([]model.User, 0)
	err := u.db.Where("address in ?", addresses).Find(&users).Error
	if err != nil {
		logger.WithError(err).Errorf("query users by addresses occur err")
		return nil, err
	}
	return users, err
}
