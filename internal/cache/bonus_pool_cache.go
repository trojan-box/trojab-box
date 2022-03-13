package cache

import (
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/app"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/jinzhu/now"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
	"time"
)

var Cache BonusPoolCache

const PoolBonusAmountKey = "PoolBonusAmountKey"

type BonusPoolCache interface {
	InitFromDB(db *gorm.DB, config app.Config)
	UpdateCache(key string, value interface{})
	GetCache(key string) interface{}
	GetTodayWonBigBonus() []constant.BonusLevel
	AddTodayWonBigBonus(bonus constant.BonusLevel)
	GetThisHourWonBigBonus() constant.BonusLevel
	SetThisHourWonBigBonus(bonus constant.BonusLevel)
	AddBonusToPool(amount int64)
	ReduceBonusFromPool(amount int64)
	GetBonusAmountFromPool() int64
}

func InitBonusPoolCache(db *gorm.DB, config app.Config) {
	Cache = &bonusPoolCache{
		cache: map[string]interface{}{},
		m:     new(sync.RWMutex),
	}
	Cache.InitFromDB(db, config)
}

type bonusPoolCache struct {
	cache map[string]interface{}
	m     *sync.RWMutex
}

func (c *bonusPoolCache) InitFromDB(db *gorm.DB, config app.Config) {

	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)

	todayBegin := now.New(shanghaiTime).BeginningOfDay().Unix()
	todayEnd := now.New(shanghaiTime).EndOfDay().Unix()

	thisHourBegin := now.New(shanghaiTime).BeginningOfHour().Unix()
	thisHourEnd := now.New(shanghaiTime).EndOfHour().Unix()

	gameRepo := repository.NewGameSessionRepo(db)
	todayWinBigBonusRecords, err := gameRepo.GetWinBigBonusByTime(todayBegin, todayEnd)
	if err != nil {
		logger.WithError(err).Errorf("get today win big bonus record occur err")
		return
	}
	for _, recordTemp := range todayWinBigBonusRecords {
		c.AddTodayWonBigBonus(recordTemp.BonusLevel)
	}

	thisHourWinBigBonusRecords, err := gameRepo.GetWinBigBonusByTime(thisHourBegin, thisHourEnd)
	if err != nil {
		logger.WithError(err).Errorf("get this hour win big bonus record occur err")
		return
	}
	for _, recordTemp := range thisHourWinBigBonusRecords {
		c.SetThisHourWonBigBonus(recordTemp.BonusLevel)
	}
	c.initPoolBonusAmount(db, config)
}

func (c *bonusPoolCache) initPoolBonusAmount(db *gorm.DB, config app.Config) {
	gameRepo := repository.NewGameSessionRepo(db)
	playGameCount, err := gameRepo.GetTotalCount()
	if err != nil {
		logger.WithError(err).Errorf("get game total count occur err")
		return
	}
	withdrawBonusRepo := repository.NewBonusWithdrawRepo(db)
	withdrawTotal, err := withdrawBonusRepo.GetTotalCount()
	if err != nil {
		logger.WithError(err).Errorf("get withdraw total count occur err")
		return
	}

	DefaultBonusPoolAmount := config.BonusPool.DefaultBonusPoolAmount
	PlayGameAddAmount := config.BonusPool.PlayGameAddAmount
	WithdrawReduceAmount := config.BonusPool.WithdrawReduceAmount
	amount := DefaultBonusPoolAmount + PlayGameAddAmount*playGameCount - withdrawTotal*WithdrawReduceAmount
	c.cache[PoolBonusAmountKey] = amount
}

func (c *bonusPoolCache) UpdateCache(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()
	c.cache[key] = value
}
func (c *bonusPoolCache) GetCache(key string) interface{} {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.cache[key]
}
func (c *bonusPoolCache) GetTodayWonBigBonus() []constant.BonusLevel {
	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)
	todayStr := shanghaiTime.Format("20060102")
	c.m.RLock()
	defer c.m.RUnlock()
	key := fmt.Sprintf("day-%s-won", todayStr)
	if value, ok := c.cache[key]; ok {
		return value.([]constant.BonusLevel)
	} else {
		return []constant.BonusLevel{}
	}
}
func (c *bonusPoolCache) AddTodayWonBigBonus(bonus constant.BonusLevel) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)
	todayStr := shanghaiTime.Format("20060102")
	c.m.Lock()
	defer c.m.Unlock()
	key := fmt.Sprintf("day-%s-won", todayStr)
	if value, ok := c.cache[key]; ok {
		bonuses := value.([]constant.BonusLevel)
		bonuses = append(bonuses, bonus)
		c.cache[key] = bonuses
	} else {
		bonuses := make([]constant.BonusLevel, 0)
		bonuses = append(bonuses, bonus)
		c.cache[key] = bonuses
	}
}
func (c *bonusPoolCache) GetThisHourWonBigBonus() constant.BonusLevel {
	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)
	thisHourStr := shanghaiTime.Format("2006010215")
	c.m.RLock()
	defer c.m.RUnlock()
	key := fmt.Sprintf("hour-%s-won", thisHourStr)
	if value, ok := c.cache[key]; ok {
		return value.(constant.BonusLevel)
	} else {
		return constant.BonusLevelUnknown
	}
}

func (c *bonusPoolCache) SetThisHourWonBigBonus(bonus constant.BonusLevel) {
	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)
	thisHourStr := shanghaiTime.Format("2006010215")
	c.m.Lock()
	defer c.m.Unlock()
	key := fmt.Sprintf("hour-%s-won", thisHourStr)
	c.cache[key] = bonus
}

func (c *bonusPoolCache) AddBonusToPool(amount int64) {
	c.m.Lock()
	defer c.m.Unlock()
	if _, ok := c.cache[PoolBonusAmountKey]; ok {
		curAmount := c.cache[PoolBonusAmountKey].(int64)
		c.cache[PoolBonusAmountKey] = curAmount + amount
	} else {
		c.cache[PoolBonusAmountKey] = amount
	}
}

func (c *bonusPoolCache) ReduceBonusFromPool(amount int64) {
	c.m.Lock()
	defer c.m.Unlock()
	if _, ok := c.cache[PoolBonusAmountKey]; ok {
		curAmount := c.cache[PoolBonusAmountKey].(int64)
		c.cache[PoolBonusAmountKey] = curAmount - amount
	} else {
		c.cache[PoolBonusAmountKey] = 0
	}
}

func (c *bonusPoolCache) GetBonusAmountFromPool() int64 {
	c.m.RLock()
	defer c.m.RUnlock()
	if _, ok := c.cache[PoolBonusAmountKey]; ok {
		return c.cache[PoolBonusAmountKey].(int64)
	} else {
		return 0
	}
}
