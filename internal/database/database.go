// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	config "github.com/samwang0723/jarvis/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

const dsnCount = 2

func GormFactory(cfg *config.Config) *gorm.DB {
	dsns := generateDSN(cfg)
	logLevel := logger.Error
	switch cfg.Environment {
	case "local":
		logLevel = logger.Info
	case "dev":
		logLevel = logger.Info
	case "staging":
		logLevel = logger.Warn
	case "production":
		logLevel = logger.Error
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logLevel,    // Log level (Info level logs SQL statements)
			Colorful:      false,       // Enable color
		},
	)
	session, err := gorm.Open(mysql.Open(dsns["master"]), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("connect database error: " + err.Error())
	}
	if dsns["replica"] != "" {
		dbResolverCfg := dbresolver.Config{
			Replicas: []gorm.Dialector{mysql.Open(dsns["replica"])},
		}
		readWritePlugin := dbresolver.Register(dbResolverCfg)
		err = session.Use(readWritePlugin)
		if err != nil {
			panic("cannot use read/write plugin: " + err.Error())
		}
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

func generateDSN(cfg *config.Config) map[string]string {
	resp := make(map[string]string, dsnCount)
	database := cfg.Database
	masterDsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&timeout=10s",
		database.User,
		database.Password,
		fmt.Sprintf("tcp(%s:%d)", database.Host, database.Port),
		database.Database,
	)
	resp["master"] = masterDsn

	if cfg.Replica.User != "" {
		replica := cfg.Replica
		replicaDsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&timeout=10s",
			replica.User,
			replica.Password,
			fmt.Sprintf("tcp(%s:%d)", database.Host, database.Port),
			replica.Database,
		)
		resp["replica"] = replicaDsn
	}

	return resp
}

type contextTxKey struct{}

func GetTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)

	return tx, ok
}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, contextTxKey{}, tx)
}

type DBTX interface {
	Omit(columns ...string) (tx *gorm.DB)
	Save(value interface{}) (tx *gorm.DB)
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	WithContext(ctx context.Context) *gorm.DB
}

type Query struct {
	txdb DBTX
}

func NewQuery(db *gorm.DB) *Query {
	return &Query{txdb: db}
}

func (q *Query) WithTx(tx *gorm.DB) *Query {
	return &Query{txdb: tx}
}

func (q *Query) Save(value interface{}) *gorm.DB {
	return q.txdb.Save(value)
}

func (q *Query) Where(query interface{}, args ...interface{}) *gorm.DB {
	return q.txdb.Where(query, args...)
}

func (q *Query) WithContext(ctx context.Context) *gorm.DB {
	return q.txdb.WithContext(ctx)
}

func (q *Query) Omit(columns ...string) *gorm.DB {
	return q.txdb.Omit(columns...)
}
