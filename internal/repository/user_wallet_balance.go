package repository

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserWalletBalanceRepo interface {
	Save(userWalletBalance model.UserWalletBalance) error
	GetTotalAmountByTime(startTime, endTime int64) (int64, error)
	GetTotalAmount() (int64, error)
}

func NewUserWalletBalance(db *gorm.DB) UserWalletBalanceRepo {
	return &userWalletBalanceRepo{db}
}

type userWalletBalanceRepo struct {
	db *gorm.DB
}

func (u *userWalletBalanceRepo) Save(userWalletBalance model.UserWalletBalance) error {
	log := logger.WithField("address", userWalletBalance.Address)
	err := u.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}, {Name: "time"}},
		DoNothing: true,
	}).Create(&userWalletBalance).Error
	if err != nil {
		log.WithError(err).Errorf("save userWalletBalance occur error")
		return err
	}
	return nil
}

func (r *userWalletBalanceRepo) GetTotalAmountByTime(startTime, endTime int64) (int64, error) {
	var sum int64
	err := r.db.Raw("select IFNULL(sum(a.balance), 0) from t_user_wallet_balance a,(select address,max(time) "+
		"as time from t_user_wallet_balance where created_at>=? and created_at<=? group by address ) as b "+
		"where a.address = b.address and a.time = b.time;", startTime, endTime).Scan(&sum).Error
	if err != nil {
		logger.WithError(err).Error("get sum balance by time occur error")
		return 0, err
	}
	return sum, nil
}
func (r *userWalletBalanceRepo) GetTotalAmount() (int64, error) {
	var sum int64
	err := r.db.Raw("select IFNULL(sum(a.balance), 0) from t_user_wallet_balance a,(select address,max(time) " +
		"as time from t_user_wallet_balance group by address ) as b " +
		"where a.address = b.address and a.time = b.time;").Scan(&sum).Error
	if err != nil {
		logger.WithError(err).Error("get sum balance by time occur error")
		return 0, err
	}
	return sum, nil
}
