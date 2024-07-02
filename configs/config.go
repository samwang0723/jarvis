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
package config

import (
	"fmt"
	"os"

	"github.com/samwang0723/jarvis/internal/helper"
	"gopkg.in/yaml.v3"
)

const (
	SecretUsername = "SECRET_USERNAME"
	SecretPassword = "SECRET_PASSWORD"
	RedisPassword  = "REDIS_PASSWD"
	SentryDSN      = "SENTRY_DSN"
)

type Config struct {
	RedisCache struct {
		Master        string   `yaml:"master"`
		Password      string   `yaml:"password"`
		SentinelAddrs []string `yaml:"sentinelAddrs"`
	} `yaml:"redis"`
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	Environment string
	Kafka       struct {
		GroupID string   `yaml:"groupId"`
		Brokers []string `yaml:"brokers"`
		Topics  []string `yaml:"topics"`
	} `yaml:"kafka"`
	Sentry struct {
		DSN   string `yaml:"dsn"`
		Debug bool   `yaml:"debug"`
	} `yaml:"sentry"`
	Server struct {
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Version  string `yaml:"version"`
		Port     int    `yaml:"port"`
		GrpcPort int    `yaml:"grpcPort"`
	} `yaml:"server"`
	ElasticSearch struct {
		Host                string `yaml:"host"`
		Port                int    `yaml:"port"`
		HealthCheckInterval int    `yaml:"healthCheckInterval"`
	} `yaml:"elasticsearch"`
	Database struct {
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Host         string `yaml:"host"`
		Database     string `yaml:"database"`
		Port         int    `yaml:"port"`
		MaxLifetime  int    `yaml:"maxLifetime"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
	} `yaml:"database"`
	Replica struct {
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Host         string `yaml:"host"`
		Database     string `yaml:"database"`
		Port         int    `yaml:"port"`
		MaxLifetime  int    `yaml:"maxLifetime"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
	} `yaml:"replica"`
}

//nolint:nolintlint, gochecknoglobals
var instance Config

func Load() {
	env := helper.GetCurrentEnv()
	yamlFile := fmt.Sprintf("./configs/config.%s.yaml", env)

	configFile, err := os.Open(yamlFile)
	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&instance)
	if err != nil {
		panic(err)
	}

	if user := os.Getenv(SecretUsername); user != "" {
		instance.Database.User = user
		instance.Replica.User = user
	}

	if passwd := os.Getenv(SecretPassword); passwd != "" {
		instance.Database.Password = passwd
		instance.Replica.Password = passwd
	}

	if redisPasswd := os.Getenv(RedisPassword); redisPasswd != "" {
		instance.RedisCache.Password = redisPasswd
	}

	if dsn := os.Getenv(SentryDSN); dsn != "" {
		instance.Sentry.DSN = dsn
	}

	instance.Environment = env
}

func GetCurrentConfig() *Config {
	return &instance
}
