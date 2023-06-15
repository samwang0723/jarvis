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
package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	redis "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
)

const (
	CronjobStockListLock = "jarvis-stock-list-lock"
	CronjobLock          = "jarvis-realtime-lock"
)

//go:generate mockgen -source=redis.go -destination=mocks/redis.go -package=cache
type Redis interface {
	SetExpire(ctx context.Context, key string, expired time.Time) error
	SAdd(ctx context.Context, key string, values []string) error
	SMembers(ctx context.Context, key string) ([]string, error)
	Set(ctx context.Context, key, val string, expired time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	MGet(ctx context.Context, keys ...string) ([]string, error)
	ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock
	Close() error
}

// Config encapsulates the settings for configuring the redis service.
type Config struct {
	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *zerolog.Logger

	// Redis master node DNS hostname
	Master string

	// Redis sentinel addresses
	SentinelAddrs []string

	// Redis password
	Password string
}

type redisImpl struct {
	instance *redis.Client
	cfg      Config
}

func New(cfg Config) Redis {
	impl := &redisImpl{
		cfg: cfg,
		instance: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       cfg.Master,
			SentinelAddrs:    cfg.SentinelAddrs,
			Password:         cfg.Password,
			SentinelPassword: cfg.Password,
			DB:               0,
		}),
	}

	return impl
}

func (r *redisImpl) SetExpire(ctx context.Context, key string, expired time.Time) error {
	expire, err := r.instance.ExpireAt(ctx, key, expired).Result()
	if err != nil {
		return xerrors.Errorf("cache.SetExpire: failed, key=%s; expired=%s; err=%w;", key, expired, err)
	}

	r.cfg.Logger.Info().Msgf("cache.SetExpire: success, key=%s; expired=%t;", key, expire)

	return nil
}

func (r *redisImpl) SAdd(ctx context.Context, key string, value []string) error {
	err := r.instance.SAdd(ctx, key, value).Err()
	if err != nil {
		return xerrors.Errorf("cache.SAdd: failed, key=%s; value=%s; err=%w;", key, value, err)
	}

	r.cfg.Logger.Info().Msgf("cache.SAdd: success, key=%s; value=%s;", key, value)

	return nil
}

func (r *redisImpl) SMembers(ctx context.Context, key string) ([]string, error) {
	res, err := r.instance.SMembers(ctx, key).Result()
	if err != nil {
		return nil, xerrors.Errorf("cache.SMembers: failed, key=%s; err=%w;", key, err)
	}

	r.cfg.Logger.Info().Msgf("cache.SMembers: success, count=%d;", len(res))

	return res, nil
}

func (r *redisImpl) Set(ctx context.Context, key, val string, expired time.Duration) error {
	res, err := r.instance.Set(ctx, key, val, expired).Result()
	if err != nil {
		return xerrors.Errorf("cache.Set: failed, key=%s; err=%w;", key, err)
	}

	r.cfg.Logger.Info().Msgf("cache.Set: success, res=%+v; key=%s", res, key)

	return nil
}

func (r *redisImpl) MGet(ctx context.Context, keys ...string) ([]string, error) {
	res, err := r.instance.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, xerrors.Errorf("cache.MGet: failed, key=%v; err=%w;", keys, err)
	}

	output := make([]string, len(res))
	for i, v := range res {
		if v != nil {
			output[i] = fmt.Sprint(v)
		}
	}

	r.cfg.Logger.Info().Msgf("cache.MGet: success, count=%d;", len(output))

	return output, nil
}

func (r *redisImpl) Get(ctx context.Context, key string) (string, error) {
	res, err := r.instance.Get(ctx, key).Result()
	if err != nil {
		return "", xerrors.Errorf("cache.Get: failed, key=%s; err=%w;", key, err)
	}

	r.cfg.Logger.Info().Msgf("cache.Get: success, res=%+v;", res)

	return res, nil
}

func (r *redisImpl) ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock {
	// Create a new lock client.
	locker := redislock.New(r.instance)

	// Try to obtain lock.
	lock, err := locker.Obtain(ctx, key, expire, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		r.cfg.Logger.Error().Err(err).Msg("cache.ObtainLock: failed, could not obtain lock!")

		return nil
	} else if err != nil {
		return nil
	}

	r.cfg.Logger.Debug().Msgf("cache.ObtainLock: success, key=%s;", key)

	return lock
}

func (r *redisImpl) Close() error {
	if err := r.instance.Close(); err != nil {
		return xerrors.Errorf("cache.Close: failed, err=%w;", err)
	}

	return nil
}
