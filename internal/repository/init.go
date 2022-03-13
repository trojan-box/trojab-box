package repository

import (
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/model"
	"github.com/aresprotocols/trojan-box/internal/pkg/config"
	"github.com/aresprotocols/trojan-box/internal/pkg/storage/orm"
	logger "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// DB Global
var DB *gorm.DB

// Init init databse
func Init() *gorm.DB {
	cfg, err := loadConf()
	if err != nil {
		panic(fmt.Sprintf("load orm conf err: %v", err))
	}

	DB = orm.NewMySQL(cfg)
	initTable(DB)
	return DB
}

// GetDB return default database
func GetDB() *gorm.DB {
	return DB
}

// loadConf load gorm config
func loadConf() (ret *orm.Config, err error) {
	var cfg orm.Config
	if err := config.Load("database", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func initTable(db *gorm.DB) {
	createTable(db, &model.Nonce{})
	createTable(db, &model.User{})
	createTable(db, &model.GameSession{})
	createTable(db, &model.UserBonus{})
	createTable(db, &model.BonusRecord{})
	createTable(db, &model.BonusWithdraw{})
	createTable(db, &model.Broadcast{})
	createTable(db, &model.Leaderboard{})
	createTable(db, &model.DailyStats{})
	createTable(db, &model.UserWalletBalance{})
	createTable(db, &model.UserYieldHourlyStats{})
	createTable(db, &model.SocialShare{})
	createTable(db, &model.UserMessage{})
	createTable(db, &model.UploadFile{})
}

func createTable(db *gorm.DB, m interface{}) {
	if db.Migrator().HasTable(m) {
		db.AutoMigrate(m)
	} else {
		logger.Debugf("create table: %v", m)
		db.Migrator().CreateTable(m)
	}
}
