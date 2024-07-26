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
	"flag"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	DbUsername    = "DB_USERNAME"
	DbPassword    = "DB_PASSWD"
	RedisPassword = "REDIS_PASSWD"
	SmartProxy    = "SMART_PROXY"
	JwtSecret     = "JWT_SECRET"
	Recaptcha     = "RECAPTCHA_SECRET"
	EnvCoreKey    = "ENVIRONMENT"
	EnvLocal      = "local"
	EnvDev        = "dev"
	EnvStaging    = "staging"
	EnvProd       = "prod"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

//nolint:gochecknoglobals // for jwt secret encryption/decryption
var JwtClaimsKey = contextKey("jwtClaims")

type Config struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	JwtSecret       string
	RecaptchaSecret string
	Kafka           struct {
		GroupID string   `yaml:"groupId"`
		Brokers []string `yaml:"brokers"`
		Topics  []string `yaml:"topics"`
	} `yaml:"kafka"`
	RedisCache struct {
		Master        string   `yaml:"master"`
		Password      string   `yaml:"password"`
		SentinelAddrs []string `yaml:"sentinelAddrs"`
	} `yaml:"redis"`
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
		Port         string `yaml:"port"`
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
		MinIdleConns int    `yaml:"minIdleConns"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
	} `yaml:"replica"`
}

//nolint:nolintlint, gochecknoglobals
var (
	instance        Config
	env             string
	dbuser          string
	dbpasswd        string
	rdpasswd        string
	jwtsecretLocal  string
	recaptchaSecret string
	migrationDown   bool
	envOnce         sync.Once
)

func GetCurrentEnv() string {
	envOnce.Do(func() {
		inputEnv := os.Getenv(EnvCoreKey)
		if inputEnv == "" {
			flag.StringVar(&inputEnv, "env", "local", "environment you want start the server")
			flag.StringVar(&dbuser, "dbuser", "", "database username")
			flag.StringVar(&jwtsecretLocal, "jwtsecret", "", "jwt secret key")
			flag.StringVar(&dbpasswd, "dbpasswd", "", "database password")
			flag.StringVar(&rdpasswd, "rdpasswd", "", "redis password")
			flag.BoolVar(&migrationDown, "down", false, "migrate down")
			flag.StringVar(&recaptchaSecret, "recaptcha", "", "reCaptcha secret")

			flag.Parse()
		}
		switch inputEnv {
		case "dev":
			env = EnvDev
		case "staging":
			env = EnvStaging
		case "prod":
			env = EnvProd
		default:
			env = EnvLocal
		}
	})

	return env
}

func IsMigrationDown() bool {
	return migrationDown
}

func Load() {
	env := GetCurrentEnv()
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

	if user := os.Getenv(DbUsername); user != "" {
		instance.Database.User = user
		instance.Replica.User = user
	} else {
		instance.Database.User = dbuser
		instance.Replica.User = dbuser
	}

	if passwd := os.Getenv(DbPassword); passwd != "" {
		instance.Database.Password = passwd
		instance.Replica.Password = passwd
	} else {
		instance.Database.Password = dbpasswd
		instance.Replica.Password = dbpasswd
	}

	if redisPasswd := os.Getenv(RedisPassword); redisPasswd != "" {
		instance.RedisCache.Password = redisPasswd
	} else {
		instance.RedisCache.Password = rdpasswd
	}

	if jwtsecret := os.Getenv(JwtSecret); jwtsecret != "" {
		instance.JwtSecret = jwtsecret
	} else {
		instance.JwtSecret = jwtsecretLocal
	}

	if recaptcha := os.Getenv(Recaptcha); recaptcha != "" {
		instance.RecaptchaSecret = recaptcha
	} else {
		instance.RecaptchaSecret = recaptchaSecret
	}
}

func GetCurrentConfig() *Config {
	return &instance
}
