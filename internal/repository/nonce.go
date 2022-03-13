package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type NonceRepo interface {
	InsertNonce(address string, nonce string) (int64, error)
	FindByAddressAndNonce(address string, nonce string) (*model.Nonce, error)
	RemoveNonce(address string) error
}

func NewNonceRepo(db *gorm.DB) NonceRepo {
	return &nonceRepo{db}
}

type nonceRepo struct {
	db *gorm.DB
}

func (r *nonceRepo) InsertNonce(address string, nonce string) (int64, error) {

	nonceModel := model.Nonce{
		Address: address,
		Nonce:   nonce,
	}
	err := r.db.Model(&nonceModel).Create(&nonceModel).Error
	if err != nil {
		logger.WithError(err).Errorf("insert nonce occur error:address:%s,nonce:%s", address, nonce)
		return 0, err
	}
	return nonceModel.ID, nil
}

func (r *nonceRepo) FindByAddressAndNonce(address string, nonce string) (*model.Nonce, error) {
	nonceModal := &model.Nonce{}
	err := r.db.Model(nonceModal).Where("address = ? and nonce = ?", address, nonce).First(nonceModal).Error
	if err != nil {
		logger.WithError(err).Errorf("find nonce by address and nonce occur error,address:%s,nonce:%s", address, nonce)
		return nil, err
	}
	return nonceModal, nil
}

func (r *nonceRepo) RemoveNonce(address string) error {

	ids := make([]int64, 0)
	err := r.db.Model(&model.Nonce{}).Select("id").Where("address = ?", address).Find(&ids).Error
	if err != nil {
		logger.WithError(err).Errorf("find ids occur error,address:%s", address)
		return err
	}
	err = r.db.Model(&model.Nonce{}).Unscoped().Where("id in ?", ids).Delete(&model.Nonce{}).Error
	if err != nil {
		logger.WithError(err).Errorf("delete nonce by address occur error,address:%s", address)
		return err
	}
	return nil
}
