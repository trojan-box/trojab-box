package usecase

import (
	i18n_lib "github.com/aresprotocols/trojan-box/internal/pkg/i18n"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	logger "github.com/sirupsen/logrus"
)

type BroadcastUseCase struct {
	BroadcastRepo repository.BroadcastRepo
}

func newBroadcast(svc *useCase) *BroadcastUseCase {
	return &BroadcastUseCase{
		BroadcastRepo: repository.NewBroadcast(svc.db),
	}
}

func (u *BroadcastUseCase) GetLatest(lang string) (string, error) {

	broadcast, err := u.BroadcastRepo.GetLatest()
	if err != nil {
		logger.WithError(err).Errorf("get latest from db occur error")
		return "", err
	}
	loc := i18n_lib.GetLocalizer(lang)
	message, _ := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    broadcast.TemplateKey,
		TemplateData: broadcast.TemplateData,
	})
	return message, nil
}
func (u *BroadcastUseCase) GetBroadcastsByPaging(lang string, page, size int) (int64, []string, error) {
	total, broadcasts, err := u.BroadcastRepo.PagingByAddress(page, size)
	if err != nil {
		logger.WithError(err).Errorf("get broadcasts from db occur error")
		return 0, nil, err
	}
	loc := i18n_lib.GetLocalizer(lang)
	messages := make([]string, 0)
	for _, broadcastTemp := range broadcasts {
		message, _ := loc.Localize(&i18n.LocalizeConfig{
			MessageID:    broadcastTemp.TemplateKey,
			TemplateData: broadcastTemp.TemplateData,
		})
		messages = append(messages, message)
	}
	return total, messages, nil
}
