package usecase

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/jinzhu/copier"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUseCase struct {
	userRepo repository.UserRepo
}

func newUser(svc *useCase) *UserUseCase {
	return &UserUseCase{userRepo: repository.NewUser(svc.db)}
}

func (u *UserUseCase) Auth(db *gorm.DB, authReq vo.UserAuthReq) (string, error) {

	authUseCase := Svc.Auth()

	token, err := authUseCase.GenerateToken(authReq.Address)
	if err != nil {
		logger.WithError(err).Error("generate jwt token occur error")
		return "", err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		nonceRepo := repository.NewNonceRepo(tx)
		err = nonceRepo.RemoveNonce(authReq.Address)
		if err != nil {
			logger.WithError(err).Errorf("remove nonce occur error,address:%s", authReq.Address)
			return err
		}

		userRepo := repository.NewUser(tx)
		err = userRepo.InsertUser(model.User{
			NickName: authReq.Address,
			Avatar:   0,
			Address:  authReq.Address,
		})
		if err != nil {
			logger.WithError(err).Errorf("insert user occur error,%v", authReq)
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	} else {
		return token, nil
	}

}

func (u *UserUseCase) GetUserProfile(address string) (vo.UserProfile, error) {
	log := logger.WithField("address", address)
	user, err := u.userRepo.GetUserByAddress(address)
	if err != nil {
		log.WithError(err).Errorf("get user by address occur err")
		return vo.UserProfile{}, err
	}
	profile := vo.UserProfile{}
	copier.Copy(&profile, user)
	return profile, nil
}

func (u *UserUseCase) UpdateUserProfile(user model.User) error {
	return u.userRepo.UpdateUserProfile(user)
}

func (u *UserUseCase) GetUsersByAddress(addresses []string) (map[string]model.User, error) {
	users, err := u.userRepo.GetUserByAddresses(addresses)
	if err != nil {
		logger.WithError(err).Errorf("get users by address occur error")
		return nil, err
	}
	userMap := make(map[string]model.User, 0)
	for _, u := range users {
		userMap[u.Address] = u
	}
	return userMap, nil
}
