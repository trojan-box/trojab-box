package orm

import (
	"database/sql"
	"fmt"
	"github.com/aresprotocols/trojan-box/internal/app"
	"log"
	"os"
	"time"

	// MySQL driver.
	"gorm.io/driver/mysql"
	// GORM MySQL
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config mysql config
type Config struct {
	Name            string
	Addr            string
	UserName        string
	Password        string
	ShowLog         bool
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifeTime time.Duration
	SlowThreshold   time.Duration
}

func NewMySQL(c *Config) (db *gorm.DB) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		c.UserName,
		c.Password,
		c.Addr,
		c.Name,
		true,
		//"Asia/Shanghai"),
		"Local")

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panicf("open mysql failed. database name: %s, err: %+v", c.Name, err)
	}

	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime)

	db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), gormConfig(c))
	if err != nil {
		log.Panicf("database connection failed. database name: %s, err: %+v", c.Name, err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	if app.Conf.Debug {
		db = db.Debug()
	}
	return db
}

func gormConfig(c *Config) *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	if c.ShowLog {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}
	if c.SlowThreshold > 0 {
		config.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: c.SlowThreshold, // nolint: golint
				Colorful:      true,
				LogLevel:      logger.Warn,
			},
		)
	}
	return config
}
