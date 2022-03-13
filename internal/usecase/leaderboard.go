package usecase

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/now"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type LeaderboardUseCase struct {
	leaderboardRepo repository.LeaderboardRepo
}

func newLeaderboard(svc *useCase) *LeaderboardUseCase {
	return &LeaderboardUseCase{
		leaderboardRepo: repository.NewLeaderboard(svc.db),
	}
}

func (u *LeaderboardUseCase) GenerateNewStar(db *gorm.DB) {
	logger.Info("start generate new start leaderboard")
	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)
	yesterday := shanghaiTime.Add(time.Duration(-24) * time.Hour)
	beginTime := now.New(yesterday).BeginningOfHour()
	endTime := now.New(shanghaiTime).EndOfHour()
	beginTimestamp := beginTime.Unix()
	endTimestamp := endTime.Unix()

	err := db.Transaction(func(tx *gorm.DB) error {
		gameSessionRepo := repository.NewGameSessionRepo(tx)
		newStartRecord, err := gameSessionRepo.GetRankByReward(beginTimestamp, endTimestamp, 3)
		if err != nil {
			logger.WithError(err).Errorf("get rank by total win occur error")
			return err
		}
		leaderboards := make([]model.Leaderboard, 0)
		for _, newStarTemp := range newStartRecord {
			leaderboard := model.Leaderboard{
				Address: newStarTemp.Address,
				Reward:  newStarTemp.Reward,
				Type:    model.LeaderboardTypeNewStar,
			}
			leaderboards = append(leaderboards, leaderboard)
		}

		leaderboardRepo := repository.NewLeaderboard(tx)
		err = leaderboardRepo.DeleteByType(model.LeaderboardTypeNewStar)
		if err != nil {
			logger.WithError(err).Errorf("delete old new start occur error")
			return err
		}
		_, err = leaderboardRepo.Save(leaderboards)
		if err != nil {
			logger.WithError(err).Errorf("save new new start occur error")
			return err
		}
		return nil
	})
	if err != nil {
		logger.WithError(err).Errorf("generate new start error")
	}

	return
}

func (u *LeaderboardUseCase) GenerateSeasonChampion(db *gorm.DB) {
	logger.Info("start generate season champion leaderboard")
	err := db.Transaction(func(tx *gorm.DB) error {
		userBonusRepo := repository.NewUserBonusRepo(tx)
		userBonus, err := userBonusRepo.GetRankByTotalWinDesc(10)
		if err != nil {
			logger.WithError(err).Errorf("get rank by total win occur error")
			return err
		}
		leaderboards := make([]model.Leaderboard, 0)
		for _, userBonusTemp := range userBonus {
			leaderboard := model.Leaderboard{
				Address: userBonusTemp.Address,
				Reward:  userBonusTemp.TotalWin,
				Type:    model.LeaderboardTypeSeasonChampion,
			}
			leaderboards = append(leaderboards, leaderboard)
		}

		leaderboardRepo := repository.NewLeaderboard(tx)
		err = leaderboardRepo.DeleteByType(model.LeaderboardTypeSeasonChampion)
		if err != nil {
			logger.WithError(err).Errorf("delete old season champion occur error")
			return err
		}
		_, err = leaderboardRepo.Save(leaderboards)
		if err != nil {
			logger.WithError(err).Errorf("save new session champion occur error")
			return err
		}
		return nil
	})
	if err != nil {
		logger.WithError(err).Errorf("generate session champion error")
	}
}

func (u *LeaderboardUseCase) GetLeaderboardsByType(leaderboardType model.LeaderboardType) ([]vo.LeaderboardResp, error) {
	leaderboards, err := u.leaderboardRepo.GetLeaderboardByType(leaderboardType)
	if err != nil {
		logger.WithError(err).Errorf("get leaderboards by type occur error")
		return nil, err
	}
	leaderboardResps := make([]vo.LeaderboardResp, 0)
	copier.Copy(&leaderboardResps, &leaderboards)

	addresses := make([]string, 0)
	for _, l := range leaderboards {
		addresses = append(addresses, l.Address)
	}

	userUseCase := Svc.User()
	userMap, err := userUseCase.GetUsersByAddress(addresses)
	if err != nil {
		logger.WithError(err).Errorf("get users by addresses occur err")
		return nil, err
	}
	for i, l := range leaderboardResps {
		userName := l.Address
		avatar := 0
		if v, ok := userMap[l.Address]; ok {
			userName = v.NickName
			avatar = v.Avatar
		}
		l.NickName = userName
		l.Avatar = avatar
		leaderboardResps[i] = l
	}
	return leaderboardResps, nil
}
