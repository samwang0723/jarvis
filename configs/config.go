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
package config

import (
	"fmt"
	"os"

	"github.com/samwang0723/jarvis/internal/helper"
	"gopkg.in/yaml.v2"
)

const (
	SecretUsername = "SECRET_USERNAME"
	SecretPassword = "SECRET_PASSWORD"
)

type Config struct {
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
	Server struct {
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		GrpcPort int    `yaml:"grpcPort"`
	} `yaml:"server"`
	RedisCache struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"redis"`
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	WorkerPool struct {
		MaxPoolSize  int `yaml:"maxPoolSize"`
		MaxQueueSize int `yaml:"maxQueueSize"`
	} `yaml:"workerpool"`
	ElasticSearch struct {
		Host                string `yaml:"host"`
		Port                int    `yaml:"port"`
		HealthCheckInterval int    `yaml:"healthCheckInterval"`
	} `yaml:"elasticsearch"`
}

var (
	instance Config
)

func Load() {
	yamlFile := fmt.Sprintf("./configs/config.%s.yaml", helper.GetCurrentEnv())
	f, err := os.Open(yamlFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&instance)
	if err != nil {
		panic(err)
	}

	if user := os.Getenv(SecretUsername); len(user) > 0 {
		instance.Database.User = user
		instance.Replica.User = user
	}

	if passwd := os.Getenv(SecretPassword); len(passwd) > 0 {
		instance.Database.Password = passwd
		instance.Replica.Password = passwd
	}
}

func GetCurrentConfig() *Config {
	return &instance
}
