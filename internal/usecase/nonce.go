package usecase

import (
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/repository"
	logger "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type NonceUseCase struct {
	nonceRepo repository.NonceRepo
}

func newNonce(svc *useCase) *NonceUseCase {
	return &NonceUseCase{nonceRepo: repository.NewNonceRepo(svc.db)}
}

func (u *NonceUseCase) GenerateNonce(address string) (string, error) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%08v", rnd.Int31n(100000000))
	_, err := u.nonceRepo.InsertNonce(address, vcode)
	if err != nil {
		logger.WithError(err).Errorf("insert nonce occur error,address:%s,nonce:%s", address, vcode)
		return "", err
	}
	return vcode, nil
}

func (u *NonceUseCase) FindByAddressAndNonce(address string, nonce string) (*model.Nonce, error) {
	return u.nonceRepo.FindByAddressAndNonce(address, nonce)
}
