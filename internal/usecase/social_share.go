package usecase

import (
	"errors"
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/now"
	logger "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"time"
)

type SocialShareUseCase struct {
	SocialShareRepo   repository.SocialShareRepo
	BonusWithdrawRepo repository.BonusWithdrawRepo
}

func newSocialShare(svc *useCase) *SocialShareUseCase {
	return &SocialShareUseCase{
		SocialShareRepo:   repository.NewSocialShareRepo(svc.db),
		BonusWithdrawRepo: repository.NewBonusWithdrawRepo(svc.db),
	}
}

func (u *SocialShareUseCase) CreateShare(address string, req vo.AddSocialShareReq, conf app.Config) (model.SocialShare, error) {
	share := model.SocialShare{
		Address:   address,
		Link:      req.Link,
		ShareType: req.ShareType,
		Channel:   req.Channel,
		Content:   req.Content,
		Bonus:     0,
		State:     model.SocialShareStateSubmit,
	}
	if share.ShareType == model.SocialShareTypeWithdraw {
		err := u.prepProcessWithdrawShare(&share, conf)
		if err != nil {
			logger.WithError(err).Errorln("prep process withdraw share occur err")
			return model.SocialShare{}, err
		}
	}

	err := u.validateShare(share, conf)
	if err != nil {
		logger.WithError(err).Errorln("validate share occur err")
		return model.SocialShare{}, err
	}
	return u.SocialShareRepo.Save(share)
}

func (u *SocialShareUseCase) prepProcessWithdrawShare(share *model.SocialShare, conf app.Config) error {
	// query last withdraw id
	withdraw, err := u.BonusWithdrawRepo.FindLastByAddress(share.Address)
	if err != nil {
		logger.WithError(err).Errorln("find last bonus withdraw by address occur err")
		return err
	}
	share.Associate = withdraw.ID
	share.Bonus = conf.Share.WithdrawShareBonus
	return nil
}

func (u *SocialShareUseCase) validateShare(share model.SocialShare, conf app.Config) error {
	switch share.ShareType {
	case model.SocialShareTypeWithdraw:
		return u.validateWithdrawShare(share)
	case model.SocialShareTypeCommon:
		return u.validateCommonShare(share, conf)
	case model.SocialShareTypeReport:
		return nil
	}
	return errors.New("unknown share type")

}
func (u *SocialShareUseCase) validateWithdrawShare(share model.SocialShare) error {
	_, err := u.SocialShareRepo.FindByAddressAndTypeAndAssociate(share.Address, model.SocialShareTypeWithdraw, share.Associate)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			logger.WithError(err).Errorln("find withdraw share by address and associate occur err")
			return err
		}
	} else {
		logger.Errorln("the latest withdraw shared")
		return errors.New("the latest withdraw shared")
	}
}
func (u *SocialShareUseCase) validateCommonShare(share model.SocialShare, conf app.Config) error {
	// is this share channel is in allow channel
	isContains := funk.ContainsInt(conf.Share.CommonShareAllowChannel, int(share.Channel))
	if !isContains {
		logger.Errorln("not allow common channel")
		return errors.New("not allow common channel")
	}

	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)
	dayBegin := now.New(shanghaiTime).BeginningOfDay().Unix()
	dayEnd := now.New(shanghaiTime).EndOfDay().Unix()

	_, err := u.SocialShareRepo.FindByAddressAndTypeAndChannel(share.Address, model.SocialShareTypeCommon, share.Channel, dayBegin, dayEnd)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			logger.WithError(err).Errorln("find withdraw share by address and channel occur err")
			return err
		}
	} else {
		logger.Errorln("the channel shared")
		return errors.New("the channel shared")
	}
}

func (u *SocialShareUseCase) PagingAndFilter(state model.SocialShareState, address string,
	shareType model.SocialShareType, channel model.SocialShareChannel, auditorAddress string, beginTime, endTime int64, page, size int) (int64, []vo.SocialShare, error) {
	log := logger.WithFields(logger.Fields{
		"state":     state,
		"address":   address,
		"shareType": shareType,
	})

	total, socialshares, err := u.SocialShareRepo.PagingAndFilter(state, address, shareType, channel, auditorAddress, beginTime, endTime, page, size)
	if err != nil {
		log.WithError(err).Error("get social shares paging by address occur error")
		return 0, nil, err
	}
	socialsharesVo := make([]vo.SocialShare, 0)
	copier.Copy(&socialsharesVo, &socialshares)

	addresses := make([]string, 0)
	for _, g := range socialsharesVo {
		addresses = append(addresses, g.Address)
	}
	userUseCase := Svc.User()
	userMap, err := userUseCase.GetUsersByAddress(addresses)
	if err != nil {
		logger.WithError(err).Errorf("get user map occur err")
		return 0, nil, err
	}
	for i, h := range socialsharesVo {
		userName := h.Address
		if v, ok := userMap[h.Address]; ok {
			userName = v.NickName
		}
		h.NickName = userName
		socialsharesVo[i] = h
	}

	return total, socialsharesVo, nil
}

func (u *SocialShareUseCase) Process(db *gorm.DB, req vo.ProcessSocialShareReq, auditorAddress string) error {
	log := logger.WithField("id", req.ID)
	err := db.Transaction(func(tx *gorm.DB) error {
		socialShareRepo := repository.NewSocialShareRepo(tx)
		share, err := socialShareRepo.FindById(req.ID)
		if err != nil {
			log.WithError(err).Errorf("get share by id occur err")
			return err
		}
		if share.ShareType != model.SocialShareTypeWithdraw {
			share.Bonus = req.Bonus
		}
		if share.State != model.SocialShareStateSubmit {
			log.WithError(err).Errorf("share state not submit")
			return errors.New("share state not submit")
		}
		share.State = model.SocialShareStateCompleted
		share.Accept = req.Accept
		share.Reply = req.Reply
		share.Auditor = req.Auditor
		share.AuditorAddress = auditorAddress
		share, err = socialShareRepo.Save(share)
		if err != nil {
			log.WithError(err).Errorf("update share  occur err")
			return err
		}

		if req.Bonus > 0 {
			// save bonus record
			bonusRecordRepo := repository.NewBonusRecordRepo(tx)
			bonusRecord := model.BonusRecord{
				Address:   share.Address,
				Bonus:     share.Bonus,
				Type:      model.BonusRecordTypeShare,
				Associate: share.ID,
				State:     model.BonusRecordStateCompleted,
			}
			_, err = bonusRecordRepo.Save(bonusRecord)
			if err != nil {
				log.WithError(err).Errorf("save bonus record occur error")
				return err
			}
			// save or update user bonus
			userBonusRepo := repository.NewUserBonusRepo(tx)
			userBonus, err := userBonusRepo.Find(share.Address)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				log.WithError(err).Errorf("find user bonus occur error")
				return err
			}
			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				userBonus = model.UserBonus{
					Address:  share.Address,
					Balance:  share.Bonus,
					TotalWin: share.Bonus,
				}
			} else {
				userBonus.Balance += share.Bonus
				userBonus.TotalWin += share.Bonus
			}
			_, err = userBonusRepo.Save(userBonus)
			if err != nil {
				log.WithError(err).Errorf("save user bonus occur error")
				return err
			}
		}

		// save message
		if hasMessage, userMessage := u.generateMessage(share, req); hasMessage {
			userMessageRepo := repository.NewUserMessage(tx)
			_, err = userMessageRepo.Save(userMessage)
			if err != nil {
				logger.WithError(err).Errorf("save user message to db occur error")
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.WithError(err).Errorf("process social share occur err")
		return err
	}
	return nil
}

func (u *SocialShareUseCase) generateMessage(share model.SocialShare, req vo.ProcessSocialShareReq) (bool, model.UserMessage) {

	if share.ShareType == model.SocialShareTypeReport {
		userMessage := model.UserMessage{
			Address:     share.Address,
			TemplateKey: req.Reply,
			IsTemplate:  false,
			State:       model.UserMessageStateUnread,
		}
		return true, userMessage

	} else {
		userMessage := model.UserMessage{
			Address:     share.Address,
			TemplateKey: constant.ShareBonusSuccessfulMessage,
			TemplateData: model.TemplateDataMap{
				"Channel": model.SocialShareChannelName[share.Channel],
				"Bonus":   share.Bonus,
			},
			IsTemplate: true,
			State:      model.UserMessageStateUnread,
		}
		return true, userMessage
	}
}
