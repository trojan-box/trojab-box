package bonuspool

import (
	"errors"
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/cache"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/aresprotocols/trojan-box/internal/pkg/util"
	"github.com/aresprotocols/trojan-box/internal/repository"
	logger "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type PlanOne struct {
	db             *gorm.DB
	config         app.Config
	bonusPoolCache cache.BonusPoolCache
}

func NewPlanOne(db *gorm.DB, config app.Config, bonusPoolCache cache.BonusPoolCache) *PlanOne {
	return &PlanOne{
		db:             db,
		config:         config,
		bonusPoolCache: bonusPoolCache,
	}
}

func (p *PlanOne) WinBonus(address string, chosenIndex int) (int64, constant.BonusLevel, []int64) {
	winLevel := p.GetWinBonusLevel(address)
	randomBonusLevelAmounts := p.GetRandomBonusLevelAmounts()
	winBonus := randomBonusLevelAmounts[winLevel-1]
	cardsBonus := p.GenerateCardsBonus(randomBonusLevelAmounts, winBonus, chosenIndex)
	return winBonus, winLevel, cardsBonus
}

func (p *PlanOne) GenerateCardsBonus(randomBonusLevelAmounts []int64, winBonus int64, chosenIndex int) []int64 {
	shuffled := util.ShuffledSlice(randomBonusLevelAmounts)
	exchangedNums := util.ExchangeChosenIndex(shuffled, winBonus, chosenIndex)
	return exchangedNums
}

func (p *PlanOne) GetRandomBonusLevelAmounts() []int64 {
	baseBonusAmounts := p.config.BonusPool.BonusLevelAmounts
	factor := p.config.BonusPool.AnnualizedFactor
	randomBonusAmounts := util.RandomRangeSlice(baseBonusAmounts, factor, 0.1)
	return randomBonusAmounts
}

func (p *PlanOne) GetWinBonusLevel(address string) constant.BonusLevel {
	if !p.IsWinBigBonus(address) {
		if !p.IsExceedMaxWinSmallTimes(address) {
			return constant.BonusLevel9
		}
	}
	if p.IsTodayBigBonusAllHadBeenWon() {
		return constant.BonusLevel9
	}
	if p.IsThisHourBigBonusHadBeenWon() {
		return constant.BonusLevel9
	}
	return p.WinBigBonus(address)
}

func (p *PlanOne) IsExceedMaxWinSmallTimes(address string) bool {
	if p.config.Game.MaxWinSmallTimes <= 0 {
		return false
	}
	gameRepo := repository.NewGameSessionRepo(p.db)
	latestTime := int64(0)
	game, err := gameRepo.LatestWinBigBonus(address)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WithError(err).Errorf("get latest win big bouns occur err")
			return false
		}
	} else {
		latestTime = game.CreatedAt
	}
	count, err := gameRepo.CountAfterTime(address, latestTime)
	if err != nil {
		logger.WithError(err).Errorf("count after time occur err")
		return false
	}
	if count >= p.config.Game.MaxWinSmallTimes {
		return true
	}
	return false
}

func (p *PlanOne) GetWinBigBonusRate() int64 {
	return p.config.BonusPool.WinBigBonusRatio
}

func (p *PlanOne) IsWinBigBonus(address string) bool {
	winBigRatio := p.GetWinBigBonusRate()
	winBigRandomNum := util.GetRandom(0, 100)
	if winBigRandomNum <= winBigRatio {
		return true
	} else {
		return false
	}
}

func (p *PlanOne) WinBigBonus(address string) constant.BonusLevel {
	todayRemainsBigBonus := p.TodayRemainsBigBonus()
	remainsBigBonusLens := len(todayRemainsBigBonus)
	per := int64(20)
	max := per * int64(remainsBigBonusLens)
	randomNum := util.GetRandom(0, max)
	index := randomNum / per
	winBonusLevel := todayRemainsBigBonus[index]
	p.bonusPoolCache.AddTodayWonBigBonus(winBonusLevel)
	p.bonusPoolCache.SetThisHourWonBigBonus(winBonusLevel)
	return winBonusLevel
}

func (p *PlanOne) TodayRemainsBigBonus() []constant.BonusLevel {
	todayDayWonBonus := p.bonusPoolCache.GetTodayWonBigBonus()
	totalBonusLevel := constant.TotalBonusLevel
	remainsBonusLevels := funk.Subtract(totalBonusLevel, todayDayWonBonus)
	return remainsBonusLevels.([]constant.BonusLevel)
}
func (p *PlanOne) IsTodayBigBonusAllHadBeenWon() bool {
	todayRemainsBigBonus := p.TodayRemainsBigBonus()
	return len(todayRemainsBigBonus) == 0
}
func (p *PlanOne) IsThisHourBigBonusHadBeenWon() bool {
	thisHourWonBigBonus := p.bonusPoolCache.GetThisHourWonBigBonus()
	return thisHourWonBigBonus != constant.BonusLevelUnknown
}
