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

package services

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/samwang0723/jarvis/internal/cache"
	"github.com/samwang0723/jarvis/internal/cronjob"
	"github.com/samwang0723/jarvis/internal/database/dal/idal"
	"github.com/samwang0723/jarvis/internal/kafka"
)

type Option func(o *serviceImpl)

func WithDAL(dal idal.IDAL) Option {
	return func(i *serviceImpl) {
		i.dal = dal
	}
}

func WithKafka(cfg KafkaConfig) Option {
	return func(i *serviceImpl) {
		if err := cfg.validate(); err != nil {
			return
		}

		i.consumer = kafka.New(kafka.Config{
			Brokers: cfg.Brokers,
			Topics:  cfg.Topics,
			GroupID: cfg.GroupID, // having consumer group id to prevent duplication of message consumption
			Logger:  cfg.Logger,
		})
	}
}

func WithRedis(cfg RedisConfig) Option {
	return func(i *serviceImpl) {
		if err := cfg.validate(); err != nil {
			return
		}

		i.cache = cache.New(cache.Config{
			Master:        cfg.Master,
			SentinelAddrs: cfg.SentinelAddrs,
			Logger:        cfg.Logger,
			Password:      cfg.Password,
		})
	}
}

func WithCronJob(cfg CronjobConfig) Option {
	return func(i *serviceImpl) {
		i.cronjob = cronjob.New(cronjob.Config{
			Logger: cfg.Logger,
		})
	}
}

func WithLogger(logger *zerolog.Logger) Option {
	return func(i *serviceImpl) {
		i.logger = logger
	}
}

func WithProxy(client *http.Client) Option {
	return func(i *serviceImpl) {
		i.proxyClient = client
	}
}
