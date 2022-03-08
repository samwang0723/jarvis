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
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Host         string `yaml:"host"`
		Database     string `yaml:"database"`
		MaxLifetime  int    `yaml:"maxLifetime"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
	} `yaml:"database"`
	Replica struct {
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Host         string `yaml:"host"`
		Database     string `yaml:"database"`
		MaxLifetime  int    `yaml:"maxLifetime"`
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
	} `yaml:"replica"`
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	WorkerPool struct {
		MaxPoolSize  int `yaml:"maxPoolSize"`
		MaxQueueSize int `yaml:"maxQueueSize"`
	} `yaml:"workerpool"`
}

var (
	instance Config
)

func Load() {
	f, err := os.Open("config.yaml")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&instance)
	if err != nil {
		panic(err)
	}
}

func GetCurrentConfig() *Config {
	return &instance
}
