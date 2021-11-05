package db

import (
	"fmt"
	"time"

	"samwang0723/jarvis/config"
	log "samwang0723/jarvis/logger/gorm"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GormFactory(cfg *config.Config) *gorm.DB {
	dsn := generateDSN(cfg)
	session, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: log.Logger(),
	})
	if err != nil {
		panic("connect database error: " + err.Error())
	}

	sqlDB, err := session.DB()
	if err != nil {
		panic(err)
	}
	database := cfg.Database
	sqlDB.SetConnMaxLifetime(time.Duration(database.MaxLifetime) * time.Second)
	sqlDB.SetMaxOpenConns(database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(database.MaxIdleConns)

	return session
}

func generateDSN(cfg *config.Config) string {
	database := cfg.Database
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True",
		database.User,
		database.Password,
		database.Host,
		database.Database,
	)
	return dsn
}
