package usecase

import (
	"errors"
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/jinzhu/copier"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BonusWithdrawUseCase struct {
	BonusWithdrawRepo repository.BonusWithdrawRepo
}

func newBonusWithdraw(svc *useCase) *BonusWithdrawUseCase {
	return &BonusWithdrawUseCase{BonusWithdrawRepo: repository.NewBonusWithdrawRepo(svc.db)}
}

func (u *BonusWithdrawUseCase) PagingByStateAndAddress(state model.BonusWithdrawState, address string, page, size int) (int64, []vo.BonusWithdraw, error) {
	log := logger.WithField("state", state)

	total, bonusWithdraws, err := u.BonusWithdrawRepo.PagingByStateAndAddress(state, address, page, size)
	if err != nil {
		log.WithError(err).Error("get bonus withdraws paging by address occur error")
		return 0, nil, err
	}
	bonusWithdrawsVos := make([]vo.BonusWithdraw, 0)
	copier.Copy(&bonusWithdrawsVos, &bonusWithdraws)

	addresses := make([]string, 0)
	for _, g := range bonusWithdrawsVos {
		addresses = append(addresses, g.Address)
	}
	userUseCase := Svc.User()
	userMap, err := userUseCase.GetUsersByAddress(addresses)
	if err != nil {
		logger.WithError(err).Errorf("get user map occur err")
		return 0, nil, err
	}
	for i, h := range bonusWithdrawsVos {
		userName := h.Address
		if v, ok := userMap[h.Address]; ok {
			userName = v.NickName
		}
		h.NickName = userName
		bonusWithdrawsVos[i] = h
	}

	return total, bonusWithdrawsVos, nil
}

func (u *BonusWithdrawUseCase) ProcessWithdraw(db *gorm.DB, id int64, txhash string) error {
	log := logger.WithField("id", id)
	err := db.Transaction(func(tx *gorm.DB) error {
		bonusWithdrawRepo := repository.NewBonusWithdrawRepo(tx)
		withdraw, err := bonusWithdrawRepo.FindById(id)
		if err != nil {
			log.WithError(err).Errorf("find bonus withdraw by id occur err")
			return err
		}
		if withdraw.State != model.BonusWithdrawStateSubmit {
			log.Errorf("the withdraw state not submit")
			return errors.New("withdraw state not submit")
		}
		withdraw.State = model.BonusWithdrawStateCompleted
		withdraw.Txhash = txhash
		_, err = bonusWithdrawRepo.Save(withdraw)
		if err != nil {
			log.WithError(err).Errorf("update bonus withdraw state occur err")
			return err
		}
		userBonusRepo := repository.NewUserBonusRepo(tx)
		userBonus, err := userBonusRepo.Find(withdraw.Address)
		if err != nil {
			log.WithError(err).Errorf("find user bonus occur err")
			return err
		}
		userBonus.Freeze -= withdraw.Bonus
		_, err = userBonusRepo.Save(userBonus)
		if err != nil {
			log.WithError(err).Errorf("update user bonus freeze occur err")
			return err
		}
		bonusRecordRepo := repository.NewBonusRecordRepo(tx)
		bonusRecord, err := bonusRecordRepo.FindByTypeAndAssociate(model.BonusRecordTypeWithdraw, id)
		if err != nil {
			log.WithError(err).Errorf("get bonus record occur err")
			return err
		}
		bonusRecord.State = model.BonusRecordStateCompleted
		_, err = bonusRecordRepo.Save(bonusRecord)
		if err != nil {
			log.WithError(err).Errorf("update bonus record state occur err")
			return err
		}

		// save message
		userMessageRepo := repository.NewUserMessage(tx)
		userMessage := model.UserMessage{
			Address:     withdraw.Address,
			TemplateKey: constant.WithdrawBonusSuccessfulMessage,
			TemplateData: model.TemplateDataMap{
				"Bonus": withdraw.Bonus,
			},
			State: model.UserMessageStateUnread,
		}
		_, err = userMessageRepo.Save(userMessage)
		if err != nil {
			logger.WithError(err).Errorf("save user message to db occur error")
			return err
		}
		return nil
	})
	return err
}
func (u *BonusWithdrawUseCase) ReportProcessedWithdrawTxhash(db *gorm.DB, id int64, txhash string) error {
	log := logger.WithField("id", id)
	err := db.Transaction(func(tx *gorm.DB) error {
		bonusWithdrawRepo := repository.NewBonusWithdrawRepo(tx)
		withdraw, err := bonusWithdrawRepo.FindById(id)
		if err != nil {
			log.WithError(err).Errorf("find bonus withdraw by id occur err")
			return err
		}
		if withdraw.State != model.BonusWithdrawStateCompleted {
			log.Errorf("the withdraw state not completed")
			return errors.New("withdraw state not completed")
		}
		withdraw.Txhash = txhash
		_, err = bonusWithdrawRepo.Save(withdraw)
		if err != nil {
			log.WithError(err).Errorf("update bonus withdraw state occur err")
			return err
		}
		return nil
	})
	return err
}
