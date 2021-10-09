package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	User     string
	Password string
	Host     string
	Database string
}

func GormFactory(config *Config) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	dsn := generateDSN(config)
	session, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("connect database error: " + err.Error())
	}
	return session
}

func generateDSN(config *Config) string {
	dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True", config.User, config.Password, config.Host, config.Database)
	return dsn
}
