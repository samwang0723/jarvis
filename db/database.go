// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package db

import (
	"fmt"
	"time"

	"github.com/samwang0723/jarvis/config"
	log "github.com/samwang0723/jarvis/logger/gorm"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func GormFactory(cfg *config.Config) *gorm.DB {
	dsns := generateDSN(cfg)
	session, err := gorm.Open(mysql.Open(dsns["master"]), &gorm.Config{
		Logger: log.Logger(),
	})
	if err != nil {
		panic("connect database error: " + err.Error())
	}
	dbResolverCfg := dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(dsns["master"])},
		Replicas: []gorm.Dialector{mysql.Open(dsns["replica"])},
	}
	readWritePlugin := dbresolver.Register(dbResolverCfg)
	err = session.Use(readWritePlugin)
	if err != nil {
		panic("cannot use read/write plugin: " + err.Error())
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
	resp := make(map[string]string, 2)
	database := cfg.Database
	masterDsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&timeout=10s",
		database.User,
		database.Password,
		fmt.Sprintf("tcp(%s:%d)", database.Host, database.Port),
		database.Database,
	)
	resp["master"] = masterDsn
	replica := cfg.Replica
	replicaDsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&timeout=10s",
		replica.User,
		replica.Password,
		fmt.Sprintf("tcp(%s:%d)", database.Host, database.Port),
		replica.Database,
	)
	resp["replica"] = replicaDsn

	return resp
}
