package usecase

import (
	"errors"
	"github.com/aresprotocols/trojan-box/internal/model"
	i18n_lib "github.com/aresprotocols/trojan-box/internal/pkg/i18n"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	logger "github.com/sirupsen/logrus"
)

type UserMessageUseCase struct {
	UserMessageRepo repository.UserMessageRepo
}

func newUserMessage(svc *useCase) *UserMessageUseCase {
	return &UserMessageUseCase{
		UserMessageRepo: repository.NewUserMessage(svc.db),
	}
}

func (u *UserMessageUseCase) PagingByAddressAndState(lang string, address string, state model.UserMessageState, page, size int) (int64, []vo.UserMessage, error) {
	total, messages, err := u.UserMessageRepo.PagingByAddressAndState(address, state, page, size)
	if err != nil {
		logger.WithError(err).Errorf("get user message from db occur error")
		return 0, nil, err
	}
	loc := i18n_lib.GetLocalizer(lang)
	messageVos := make([]vo.UserMessage, 0)
	for _, msgTemp := range messages {
		content := ""
		if msgTemp.IsTemplate {
			content, _ = loc.Localize(&i18n.LocalizeConfig{
				MessageID:    msgTemp.TemplateKey,
				TemplateData: msgTemp.TemplateData,
			})
		} else {
			content = msgTemp.TemplateKey
		}
		messageVo := vo.UserMessage{
			ID:      msgTemp.ID,
			Content: content,
			State:   msgTemp.State,
		}
		messageVos = append(messageVos, messageVo)
	}
	return total, messageVos, nil
}

func (u *UserMessageUseCase) ReadMessage(address string, msgId int64) error {
	msg, err := u.UserMessageRepo.FindById(msgId)
	if err != nil {
		logger.WithError(err).Errorf("query message by id fron db occur err")
		return err
	}
	if msg.Address != address {
		logger.Errorf("the message not belong to you")
		return errors.New("the message not belong to you")
	}
	msg.State = model.UserMessageStateRead
	_, err = u.UserMessageRepo.Save(msg)
	if err != nil {
		logger.WithError(err).Errorf("update user message occur err")
		return err
	}
	return nil
}
