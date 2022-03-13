package usecase

import (
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/pkg/constant"
	"github.com/aresprotocols/trojan-box/internal/repository"
	"github.com/aresprotocols/trojan-box/internal/vo"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/now"
	gocache "github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	logger "github.com/sirupsen/logrus"
	"time"
)

type StatsUseCase struct {
	UserRepo                 repository.UserRepo
	GameSessionRepo          repository.GameSessionRepo
	DailyStatsRepo           repository.DailyStatsRepo
	BonusWithdrawRepo        repository.BonusWithdrawRepo
	GoCache                  *gocache.Cache
	UserWalletBalanceRepo    repository.UserWalletBalanceRepo
	UserYieldHourlyStatsRepo repository.UserYieldHourlyStatsRepo
	SocialShareUseCase       repository.SocialShareRepo
}

func newStats(svc *useCase) *StatsUseCase {
	return &StatsUseCase{
		UserRepo:                 repository.NewUser(svc.db),
		GameSessionRepo:          repository.NewGameSessionRepo(svc.db),
		DailyStatsRepo:           repository.NewDailyStats(svc.db),
		BonusWithdrawRepo:        repository.NewBonusWithdrawRepo(svc.db),
		GoCache:                  svc.goCache,
		UserWalletBalanceRepo:    repository.NewUserWalletBalance(svc.db),
		UserYieldHourlyStatsRepo: repository.NewUserYieldHourlyStats(svc.db),
		SocialShareUseCase:       repository.NewSocialShareRepo(svc.db),
	}
}

func (u *StatsUseCase) generateDaily(dayTime time.Time) {

	dayBegin := now.New(dayTime).BeginningOfDay().Unix()
	dayEnd := now.New(dayTime).EndOfDay().Unix()

	dayStr := dayTime.Format("20060102")

	newAddressAcount, err := u.UserRepo.GetNewAddressCount(dayBegin, dayEnd)
	if err != nil {
		logger.WithError(err).Errorf("get new address count occur err")
		return
	}
	partInAddressCount, err := u.GameSessionRepo.GetPartInAddressCount(dayBegin, dayEnd)
	if err != nil {
		logger.WithError(err).Errorf("get part in address count occur err")
		return
	}
	partInCount, err := u.GameSessionRepo.GetPartInCount(dayBegin, dayEnd)
	if err != nil {
		logger.WithError(err).Errorf("get part in count occur err")
		return
	}
	gameRewardAmount, err := u.GameSessionRepo.GetRewardAmount(dayBegin, dayEnd)
	if err != nil {
		logger.WithError(err).Errorf("get game reward amount occur err")
		return
	}
	bigRewardRecords, err := u.GameSessionRepo.GetWinBigBonusByTime(dayBegin, dayEnd)

	if err != nil {
		logger.WithError(err).Errorf("get win big Bonus by time occur err")
		return
	}

	withdrawAmount, err := u.BonusWithdrawRepo.GetSumByTime(dayBegin, dayEnd)
	if err != nil {
		logger.WithError(err).Errorf("get withdraw amount by time occur err")
	}

	openedBigRewardCount := len(bigRewardRecords)
	unopenedBigRewardCount := constant.TotalBigRewardCount - openedBigRewardCount

	singleMaxReward, err := u.GameSessionRepo.GetMaxMinSingleRewardSum(dayBegin, dayEnd, true)
	if err != nil {
		logger.WithError(err).Errorf("get single max reward occur err")
		return
	}
	singleMinReward, err := u.GameSessionRepo.GetMaxMinSingleRewardSum(dayBegin, dayEnd, false)
	if err != nil {
		logger.WithError(err).Errorf("get single min reward occur err")
		return
	}
	stakingAmount, err := u.UserWalletBalanceRepo.GetTotalAmountByTime(dayBegin, dayEnd)
	if err != nil {
		logger.WithError(err).Error("get staking use balance occur err")
		return
	}
	shareRewardAmount, err := u.SocialShareUseCase.GetRewardAmount(dayBegin, dayEnd)
	if err != nil {
		logger.WithError(err).Error("get share reward amount occur err")
		return
	}
	rewardAmount := gameRewardAmount + shareRewardAmount

	annualYieldRate := float64(0)
	if stakingAmount != 0 && gameRewardAmount != 0 {
		annualYieldRate, _ = decimal.NewFromInt(gameRewardAmount).Mul(decimal.NewFromInt(365)).Div(decimal.NewFromInt(stakingAmount)).Float64()
	}

	dailyStats := model.DailyStats{
		Day:               dayStr,
		NewAddress:        newAddressAcount,
		PartInAddress:     partInAddressCount,
		PartInCount:       partInCount,
		RewardAmount:      rewardAmount,
		WithdrawAmount:    withdrawAmount,
		OpenedBigReward:   openedBigRewardCount,
		UnopenedBigReward: unopenedBigRewardCount,
		SingleRewardMax:   singleMaxReward,
		SingleRewardMin:   singleMinReward,
		StakingAmount:     stakingAmount,
		GameRewardAmount:  gameRewardAmount,
		ShareRewardAmount: shareRewardAmount,
		AnnualYieldRate:   annualYieldRate,
	}
	dailyStats, err = u.DailyStatsRepo.Save(dailyStats)
	if err != nil {
		logger.WithError(err).Errorf("save daily stats occur err")
		return
	}
}

func (u *StatsUseCase) GenerateTodayDaily() {
	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)
	u.generateDaily(shanghaiTime)
}

func (u *StatsUseCase) GenerateYesterdayDaily() {
	location, _ := time.LoadLocation("Asia/Shanghai")
	shanghaiTime := time.Now().In(location)
	yesterday := shanghaiTime.Add(time.Duration(-24) * time.Hour)
	u.generateDaily(yesterday)
}

func (u *StatsUseCase) FindDailyStatsByDay(day string) (vo.DailyStats, error) {
	dailyStats, err := u.DailyStatsRepo.GetByDay(day)
	if err != nil {
		logger.WithError(err).Errorf("get daily stats occur err")
		return vo.DailyStats{}, err
	}

	dailyStatsVo := vo.DailyStats{}
	copier.Copy(&dailyStatsVo, dailyStats)
	return dailyStatsVo, err
}

func (u *StatsUseCase) FindDailyStatsByPaging(page, size int) (int64, []vo.DailyStats, error) {
	total, dailyStats, err := u.DailyStatsRepo.Paging(page, size)
	if err != nil {
		logger.WithError(err).Errorf("get daily stats by paging occur err")
		return 0, nil, err
	}
	dailyStatsVo := make([]vo.DailyStats, 0)
	copier.Copy(&dailyStatsVo, dailyStats)
	return total, dailyStatsVo, err
}
func (u *StatsUseCase) GetTotalStats() (vo.TotalStats, error) {
	totalStats, found := u.GoCache.Get("totalStats")
	if found {
		logger.Infof("get total stats from cache")
		return totalStats.(vo.TotalStats), nil
	}

	partInAddressCount, err := u.GameSessionRepo.GetPartInAddressCount(0, 0)
	if err != nil {
		logger.WithError(err).Errorf("get part in address occur err")
		return vo.TotalStats{}, err
	}

	newUsersCount, err := u.UserRepo.GetUserTotalCount()
	if err != nil {
		logger.WithError(err).Errorf("get user total count occur err")
		return vo.TotalStats{}, err
	}
	gameRewardAmount, err := u.GameSessionRepo.GetTotalRewardAmount()
	if err != nil {
		logger.WithError(err).Errorf("get game total reward occur err")
		return vo.TotalStats{}, err
	}
	totalWithdrawAmount, err := u.BonusWithdrawRepo.GetTotalWithdrawAmount()
	if err != nil {
		logger.WithError(err).Errorf("get total withdraw occur err")
		return vo.TotalStats{}, err
	}

	shareRewardAmount, err := u.SocialShareUseCase.GetTotalRewardAmount()
	if err != nil {
		logger.WithError(err).Errorf("get share total reward occur err")
		return vo.TotalStats{}, err
	}
	totalRewardAmount := gameRewardAmount + shareRewardAmount

	apyCount, apySum, err := u.DailyStatsRepo.GetSumAnnualYieldRateAndCount()
	if err != nil {
		logger.WithError(err).Errorf("get apy count ,sum from db occur err")
		return vo.TotalStats{}, err
	}
	logger.Infof("apyCount:%d, apySum:%f", apyCount, apySum)

	avgApy, _ := decimal.NewFromFloat(apySum).Div(decimal.NewFromInt(apyCount)).Float64()

	totalGamesCount, err := u.GameSessionRepo.GetTotalCount()
	if err != nil {
		logger.WithError(err).Errorf("get total game count occur err")
		return vo.TotalStats{}, err
	}

	totalStakingSum, err := u.UserWalletBalanceRepo.GetTotalAmount()
	if err != nil {
		logger.WithError(err).Errorf("get total staking amount occur err")
		return vo.TotalStats{}, err
	}
	newTotalStats := vo.TotalStats{
		SignedAddress:     newUsersCount,
		PartInAddress:     partInAddressCount,
		PartInCount:       totalGamesCount,
		RewardAmount:      totalRewardAmount,
		WithdrawAmount:    totalWithdrawAmount,
		GameRewardAmount:  gameRewardAmount,
		ShareRewardAmount: shareRewardAmount,
		StakingAmount:     totalStakingSum,
		AvgAPR:            avgApy,
	}

	u.GoCache.Set("totalStats", newTotalStats, 1*time.Minute)
	return newTotalStats, nil
}

func (u *StatsUseCase) FindUserYieldHourlyStatsByPaging(page, size int) (int64, []vo.UserYieldHourlyStats, error) {
	total, stats, err := u.UserYieldHourlyStatsRepo.Paging(page, size)
	if err != nil {
		logger.WithError(err).Errorf("get user yield hourly stats by paging occur err")
		return 0, nil, err
	}
	statsVo := make([]vo.UserYieldHourlyStats, 0)
	copier.Copy(&statsVo, stats)
	return total, statsVo, err
}
