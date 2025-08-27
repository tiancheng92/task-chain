package modal

import (
	"time"

	"github.com/tiancheng92/task-chain/internal/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxIdleConns    = 25
	maxOpenConns    = 50
	connMaxLifetime = 5 * time.Minute
)

var db *gorm.DB

func InitByDSN(dsn string) {
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	autoMigrateTables()
}

func InitByDB(gormDB *gorm.DB) {
	db = gormDB
	autoMigrateTables()
}

func autoMigrateTables() {
	err := db.AutoMigrate(
		new(TaskNode),
		new(TaskChain),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *gorm.DB {
	return db
}
