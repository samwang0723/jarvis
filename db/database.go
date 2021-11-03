package db

import (
	"fmt"
	"time"

	log "samwang0723/jarvis/logger/gorm"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	User         string
	Password     string
	Host         string
	Database     string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

func GormFactory(config *Config) *gorm.DB {
	dsn := generateDSN(config)
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
	sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	return session
}

func generateDSN(config *Config) string {
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True", config.User, config.Password, config.Host, config.Database)
	return dsn
}
