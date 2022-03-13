package usecase

import (
	"errors"
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/cache"
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/pkg/bonuspool"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/aresprotocols/trojan-box/internal/pkg/erc20"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/copier"
	"github.com/shopspring/decimal"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type GameUseCase struct {
	gameSessionRepo       repository.GameSessionRepo
	bonusPoolPlanOne      *bonuspool.PlanOne
	userWalletBalanceRepo repository.UserWalletBalanceRepo
	bonusPoolCache        cache.BonusPoolCache
}

func newGame(svc *useCase) *GameUseCase {
	return &GameUseCase{
		gameSessionRepo:       repository.NewGameSessionRepo(svc.db),
		bonusPoolPlanOne:      bonuspool.NewPlanOne(svc.db, *app.Conf, svc.boolCache),
		userWalletBalanceRepo: repository.NewUserWalletBalance(svc.db),
		bonusPoolCache:        svc.boolCache,
	}
}

func (u *GameUseCase) VerifyWalletBalance(address string, appConfig app.Config) (int64, error) {
	aresConfig := appConfig.Ares
	client, err := ethclient.Dial(aresConfig.ApiUrl)
	if err != nil {
		logger.WithError(err).Errorf("Failed to connect to the ethereum client")
		return 0, err
	}
	ens, err := erc20.NewToken(common.HexToAddress(aresConfig.AresContractAddress), client)
	if err != nil {
		logger.WithError(err).Errorf("new token instance occur error")
		return 0, err
	}
	balance, err := ens.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		logger.WithError(err).Errorf("query addr:%s balance occur error", address)
		return 0, err
	}
	minBalance := appConfig.Game.MinBalance
	aresDecimals := aresConfig.AresContractDecimals
	minBalanceDec := decimal.New(int64(minBalance), int32(aresDecimals))
	if decimal.NewFromBigInt(balance, 0).LessThan(minBalanceDec) {
		logger.Errorf("balance less then minimum，balance：%s,minimum:%s", balance.String(), minBalanceDec.String())
		return 0, errors.New("balance less then minimum")
	}
	balanceNumber := decimal.NewFromBigInt(balance, 0).Div(decimal.New(1, int32(aresDecimals))).IntPart()
	return balanceNumber, nil
}
func (u *GameUseCase) VerifyInTimeRange(gameConfig app.GameConfig) error {
	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)
	timeHour := shanghaiTime.Hour()
	if timeHour < gameConfig.StartHour || timeHour > gameConfig.EndHour {
		logger.Errorf("time not in hour 12-24,current is:%d", timeHour)
		return errors.New(fmt.Sprintf("time not in hour 12-24"))
	}
	return nil
}

func (u *GameUseCase) VerifyChosen(chosenIndex int) error {
	if chosenIndex < 0 || chosenIndex > 8 {
		return errors.New("incorrect chosen index")
	}
	return nil
}

func (u *GameUseCase) VerifyCards(cards []constant.God) error {
	if len(cards) != 9 {
		return errors.New("cards length not equal 9")
	}

	num := make(map[constant.God]bool)
	for _, god := range cards {
		if god > 9 || god < 1 {
			return errors.New("unknown god")
		}
		if !num[god] {
			num[god] = true
		} else {
			return errors.New("repeat god")
		}
	}
	return nil
}
func (u *GameUseCase) VerifyPlayedGame(address string, sessionStr string) error {
	_, err := u.gameSessionRepo.FindByAddressAndSession(address, sessionStr)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return errors.New("played game this hour")
	}

}

func (u *GameUseCase) PlayGame(db *gorm.DB, req vo.PlayGameReq, config app.Config) (*model.GameSession, error) {
	log := logger.WithFields(logger.Fields{
		"address": req.Address,
	})
	err := u.VerifyInTimeRange(config.Game)
	if err != nil {
		log.WithError(err).Errorf("verify in time fail")
		return nil, err
	}
	userWalletBalanceNumber := int64(0)
	if config.Game.VerifyWalletBalance {
		userWalletBalanceNumber, err = u.VerifyWalletBalance(req.Address, config)
		if err != nil {
			log.WithError(err).Error("verifyWalletBalance occur err")
			return nil, err
		}
	}
	err = u.VerifyCards(req.Cards)
	if err != nil {
		log.WithError(err).Errorf("verify cards fail")
		return nil, err
	}
	err = u.VerifyChosen(req.ChosenIndex)
	if err != nil {
		log.WithError(err).Errorf("verify chosen fail")
		return nil, err
	}
	sessionStr := u.GetGameSessionStr(config.Game)
	err = u.VerifyPlayedGame(req.Address, sessionStr)
	if err != nil {
		log.WithError(err).Errorf("pladyed this hour")
		return nil, err
	}
	winBonus, bonusLevel, cardsBonus := u.bonusPoolPlanOne.WinBonus(req.Address, req.ChosenIndex)
	gameSession := model.GameSession{
		Address:     req.Address,
		Session:     sessionStr,
		ChosenIndex: req.ChosenIndex,
		Bonus:       winBonus,
		BonusLevel:  bonusLevel,
		Cards:       req.Cards,
		CardsBonus:  cardsBonus,
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		// save game session
		gameSessionRepo := repository.NewGameSessionRepo(tx)
		gameSession, err := gameSessionRepo.Save(gameSession)
		if err != nil {
			logger.WithError(err).Errorf("save game session occur error")
			return err
		}
		// save bonus record
		bonusRecordRepo := repository.NewBonusRecordRepo(tx)
		bonusRecord := model.BonusRecord{
			Address:   req.Address,
			Bonus:     winBonus,
			Type:      model.BonusRecordTypeWin,
			Associate: gameSession.ID,
			State:     model.BonusRecordStateCompleted,
		}
		_, err = bonusRecordRepo.Save(bonusRecord)
		if err != nil {
			logger.WithError(err).Errorf("save bonus record occur error")
			return err
		}
		// save or update user bonus
		userBonusRepo := repository.NewUserBonusRepo(tx)
		userBonus, err := userBonusRepo.Find(req.Address)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WithError(err).Errorf("find user bonus occur error")
			return err
		}
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			userBonus = model.UserBonus{
				Address:  req.Address,
				Balance:  winBonus,
				TotalWin: winBonus,
			}
		} else {
			userBonus.Balance += winBonus
			userBonus.TotalWin += winBonus
		}
		_, err = userBonusRepo.Save(userBonus)
		if err != nil {
			logger.WithError(err).Errorf("save user bonus occur error")
			return err
		}
		//remove nonce
		nonceRepo := repository.NewNonceRepo(tx)
		err = nonceRepo.RemoveNonce(req.Address)
		if err != nil {
			logger.WithError(err).Errorf("remove nonce occur error")
			return err
		}
		//add broadcast
		if winBonus >= config.Game.MinBroadcastBonus {
			err = u.saveBroadcast(tx, req, winBonus)
			if err != nil {
				log.WithError(err).Errorf("save broadcast occur error")
				return err
			}
		}

		location, _ := time.LoadLocation("Asia/Shanghai")
		shanghaiTime := time.Now().In(location)
		hourTimeStr := shanghaiTime.Format("2006010215")

		userWalletBalance := model.UserWalletBalance{
			Address: req.Address,
			Time:    hourTimeStr,
			Balance: userWalletBalanceNumber,
		}
		err = u.userWalletBalanceRepo.Save(userWalletBalance)
		if err != nil {
			logger.Errorf("save user wallet balance occur error")
			return err
		}

		u.bonusPoolCache.AddBonusToPool(config.BonusPool.PlayGameAddAmount)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &gameSession, nil
}

func (u *GameUseCase) saveBroadcast(tx *gorm.DB, req vo.PlayGameReq, winBonus int64) error {
	userRepo := repository.NewUser(tx)
	user, err := userRepo.GetUserByAddress(req.Address)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = model.User{
				Address:  req.Address,
				NickName: req.Address,
				Avatar:   0,
			}
		} else {
			logger.WithError(err).Errorf("query user by address occur error")
			return err
		}
	}
	broadcastRepo := repository.NewBroadcast(tx)
	broadcast := model.Broadcast{
		Address:     req.Address,
		TemplateKey: constant.WIN_CONGRATULATION,
		TemplateData: model.TemplateDataMap{
			"Name":  user.NickName,
			"Bonus": winBonus,
		},
	}
	_, err = broadcastRepo.Save(broadcast)
	if err != nil {
		logger.WithError(err).Errorf("save braodcast to db occur error")
		return err
	}
	return nil
}

func (u *GameUseCase) GenerateCardBonus() []int {
	bonusArr := make([]int, 9)
	for i := 0; i < 9; i++ {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		bonus := rnd.Int31n(1000)
		bonusArr[i] = int(bonus)
	}
	return bonusArr
}

func (u *GameUseCase) GetGameSessionStr(config app.GameConfig) string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)
	return shanghaiTime.Format(config.SessionLayout)
}

func (u *GameUseCase) GetGameDetail(address string, session string) (vo.GameSession, error) {
	log := logger.WithField("address", address).WithField("session", session)

	gameSessionModel, err := u.gameSessionRepo.FindByAddressAndSession(address, session)
	if err != nil {
		log.WithError(err).Errorf("find game session by addresss and session occur error")
		return vo.GameSession{}, err
	}
	gameSession := vo.GameSession{}
	copier.Copy(&gameSession, gameSessionModel)
	return gameSession, nil
}
func (u *GameUseCase) GetGameDetailById(address string, id int64) (vo.GameSession, error) {
	log := logger.WithField("address", address).WithField("id", id)

	gameSessionModel, err := u.gameSessionRepo.FindByAddressAndId(address, id)
	if err != nil {
		log.WithError(err).Errorf("find game session by addresss and id occur error")
		return vo.GameSession{}, err
	}
	gameSession := vo.GameSession{}
	copier.Copy(&gameSession, gameSessionModel)
	return gameSession, nil
}

func (u *GameUseCase) GetMyHistoryByPage(address string, page, size int) (int64, []vo.GameHistory, error) {
	log := logger.WithField("address", address)

	total, gameModels, err := u.gameSessionRepo.PagingByAddress(address, page, size)
	if err != nil {
		log.WithError(err).Error("get session pading by address occur error")
		return 0, nil, err
	}
	histories := make([]vo.GameHistory, 0)
	copier.Copy(&histories, &gameModels)
	return total, histories, nil
}
func (u *GameUseCase) GetHistoryByPageAndAddress(address string, page, size int) (int64, []vo.GameHistories, error) {

	total, gameModels, err := u.gameSessionRepo.PagingByAddress(address, page, size)
	if err != nil {
		logger.WithError(err).Error("get game paging occur error")
		return 0, nil, err
	}
	histories := make([]vo.GameHistories, 0)
	copier.Copy(&histories, &gameModels)

	addresses := make([]string, 0)
	for _, g := range gameModels {
		addresses = append(addresses, g.Address)
	}
	userUseCase := Svc.User()
	userMap, err := userUseCase.GetUsersByAddress(addresses)
	if err != nil {
		logger.WithError(err).Errorf("get user map occur err")
		return 0, nil, err
	}
	for i, h := range histories {
		userName := h.Address
		if v, ok := userMap[h.Address]; ok {
			userName = v.NickName
		}
		h.NickName = userName
		histories[i] = h
	}

	return total, histories, nil
}
