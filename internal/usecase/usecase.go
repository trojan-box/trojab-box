package usecase

import (
	"github.com/aresprotocols/trojan-box/internal/cache"
	gocache "github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

// Svc global var
var Svc UseCase

type UseCase interface {
	Auth() *AuthUseCase
	Nonce() *NonceUseCase
	User() *UserUseCase
	Game() *GameUseCase
	Bonus() *BonusUseCase
	Broadcast() *BroadcastUseCase
	Leaderboard() *LeaderboardUseCase
	Stats() *StatsUseCase
	BonusWithdraw() *BonusWithdrawUseCase
	PoolCache() cache.BonusPoolCache
	SocialShare() *SocialShareUseCase
	UserMessage() *UserMessageUseCase
	File() *FileUseCase
	Gas() *GasUseCase
}

type useCase struct {
	db        *gorm.DB
	boolCache cache.BonusPoolCache
	goCache   *gocache.Cache
}

func New(db *gorm.DB, boolCache cache.BonusPoolCache, goCache *gocache.Cache) UseCase {
	return &useCase{
		db:        db,
		boolCache: boolCache,
		goCache:   goCache,
	}
}

func (u *useCase) Auth() *AuthUseCase {
	return newAuth()
}

func (u *useCase) Nonce() *NonceUseCase {
	return newNonce(u)
}

func (u *useCase) User() *UserUseCase {
	return newUser(u)
}

func (u *useCase) Game() *GameUseCase {
	return newGame(u)
}
func (u *useCase) Bonus() *BonusUseCase {
	return newBonus(u)
}
func (u *useCase) Broadcast() *BroadcastUseCase {
	return newBroadcast(u)
}

func (u *useCase) Leaderboard() *LeaderboardUseCase {
	return newLeaderboard(u)
}

func (u *useCase) Stats() *StatsUseCase {
	return newStats(u)
}

func (u *useCase) BonusWithdraw() *BonusWithdrawUseCase {
	return newBonusWithdraw(u)
}

func (u *useCase) PoolCache() cache.BonusPoolCache {
	return u.boolCache
}

func (u *useCase) SocialShare() *SocialShareUseCase {
	return newSocialShare(u)
}

func (u *useCase) UserMessage() *UserMessageUseCase {
	return newUserMessage(u)
}

func (u *useCase) File() *FileUseCase {
	return newFile(u)
}

func (u *useCase) Gas() *GasUseCase {
	return newGas(u)
}
