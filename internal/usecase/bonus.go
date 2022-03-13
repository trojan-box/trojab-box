package usecase

import (
	"errors"
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/cache"
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/now"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type BonusUseCase struct {
	UserBonusRepo   repository.UserBonusRepo
	BonusRecordRepo repository.BonusRecordRepo
	bonusPoolCache  cache.BonusPoolCache
}

func newBonus(svc *useCase) *BonusUseCase {
	return &BonusUseCase{
		UserBonusRepo:   repository.NewUserBonusRepo(svc.db),
		BonusRecordRepo: repository.NewBonusRecordRepo(svc.db),
		bonusPoolCache:  svc.boolCache,
	}
}

func (u *BonusUseCase) GetUserBonus(address string) (vo.UserBonus, error) {
	log := logger.WithField("address", address)
	userBonus, err := u.UserBonusRepo.Find(address)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Infof("record not found")
			userBonus = model.UserBonus{
				Address:  address,
				Balance:  0,
				TotalWin: 0,
				Freeze:   0,
			}
		} else {
			log.WithError(err).Errorf("find user bonus by address occur error")
			return vo.UserBonus{}, err
		}
	}
	myBonus := vo.UserBonus{}
	copier.Copy(&myBonus, &userBonus)

	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)

	todayBegin := now.New(shanghaiTime).BeginningOfDay().Unix()
	todayEnd := now.New(shanghaiTime).EndOfDay().Unix()

	todayWin := int64(0)
	todayBonusRecords, err := u.BonusRecordRepo.GetByTimeAndAddress(todayBegin, todayEnd, address)

	if err != nil {
		log.WithError(err).Errorf("find today bonus record occur err")
		return vo.UserBonus{}, err
	}
	for _, r := range todayBonusRecords {
		if r.Type == model.BonusRecordTypeWin {
			todayWin += r.Bonus
		}
	}
	myBonus.TodayWin = todayWin
	return myBonus, nil
}

func (u *BonusUseCase) GetBonusHistoryByPage(address string, page, size int, bonusRecordType int) (int64, []vo.BonusHistory, error) {
	log := logger.WithField("address", address)

	total, records, err := u.BonusRecordRepo.PagingByAddress(address, page, size, bonusRecordType)
	if err != nil {
		log.WithError(err).Error("get bonus record paging by address occur error")
		return 0, nil, err
	}
	histories := make([]vo.BonusHistory, 0)
	copier.Copy(&histories, &records)
	return total, histories, nil
}

func (u *BonusUseCase) WithdrawBonusApply(db *gorm.DB, address string, withdrawBonus int64, config app.Config) error {
	log := logger.WithField("address", address)

	userBonus, err := u.UserBonusRepo.Find(address)
	if err != nil {
		log.WithError(err).Errorf("find user bonus by address occur err")
		return err
	}
	// verify balance over then 1000
	if withdrawBonus < config.Game.MinWithdraw {
		log.Errorf("withdrawBonus less then MinWithdraw")
		return errors.New("withdrawBonus less then MinWithdraw")
	}
	// verify balance over withdraw bonus
	if userBonus.Balance < withdrawBonus {
		log.Errorf("balance less then withdrawBonus:%d", withdrawBonus)
		return errors.New("balance less then withdrawBonus")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// update user bonus
		userBonusRepo := repository.NewUserBonusRepo(tx)
		userBonus.Balance -= withdrawBonus
		userBonus.Freeze += withdrawBonus
		_, err := userBonusRepo.Save(userBonus)
		if err != nil {
			logger.WithError(err).Errorf("update user bonus occur error")
			return err
		}
		// save withdraw apply
		bonusWithdrawRepo := repository.NewBonusWithdrawRepo(tx)
		bonusWithdraw := model.BonusWithdraw{
			Address: address,
			Bonus:   withdrawBonus,
			State:   model.BonusWithdrawStateSubmit,
		}
		bonusWithdraw, err = bonusWithdrawRepo.Save(bonusWithdraw)
		if err != nil {
			logger.WithError(err).Errorf("save bonus withdraw occur error")
			return err
		}
		//save bonus record
		bonusRecordRepo := repository.NewBonusRecordRepo(tx)
		bonusRecord := model.BonusRecord{
			Address:   address,
			Bonus:     withdrawBonus,
			Type:      model.BonusRecordTypeWithdraw,
			Associate: bonusWithdraw.ID,
			State:     model.BonusRecordStateSubmit,
		}
		_, err = bonusRecordRepo.Save(bonusRecord)
		if err != nil {
			logger.WithError(err).Errorf("save bonus record occur error")
			return err
		}
		//remove nonce
		nonceRepo := repository.NewNonceRepo(tx)
		err = nonceRepo.RemoveNonce(address)
		if err != nil {
			logger.WithError(err).Errorf("remove nonce occur error")
			return err
		}
		u.bonusPoolCache.ReduceBonusFromPool(config.BonusPool.WithdrawReduceAmount)
		return nil
	})

	if err != nil {
		log.WithError(err).Errorf("withdraw bonus save db occur err")
		return err
	}
	return nil
}
