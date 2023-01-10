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
	"github.com/samwang0723/jarvis/internal/cache"
	"github.com/samwang0723/jarvis/internal/cronjob"
	"github.com/samwang0723/jarvis/internal/db/dal/idal"
	"github.com/samwang0723/jarvis/internal/kafka/ikafka"
)

type Option func(o *serviceImpl)

func WithDAL(dal idal.IDAL) Option {
	return func(i *serviceImpl) {
		i.dal = dal
	}
}

func WithKafka(consumer ikafka.IKafka) Option {
	return func(i *serviceImpl) {
		i.consumer = consumer
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
