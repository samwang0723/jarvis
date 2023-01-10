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
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
)

const (
	realTimeMonitoringKey = "real_time_monitoring_keys"
	defaultCacheExpire    = 7 * 24 * time.Hour
)

func (s *serviceImpl) ListSelections(ctx context.Context,
	req *dto.ListSelectionRequest,
) ([]*entity.Selection, error) {
	objs, err := s.dal.ListSelections(ctx, req.Date, req.Strict)
	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (s *serviceImpl) PresetRealTimeKeys(ctx context.Context) error {
	keys, err := s.dal.GetRealTimeMonitoringKeys(ctx)
	if err != nil {
		return err
	}

	t := time.Now().AddDate(0, 0, 0)
	date := t.Format("20060102")
	redisKey := fmt.Sprintf("%s:%s", realTimeMonitoringKey, date)

	err = s.cache.SAdd(ctx, redisKey, keys)
	if err != nil {
		return err
	}

	err = s.cache.SetExpire(ctx, redisKey, time.Now().Add(defaultCacheExpire))
	if err != nil {
		return err
	}

	return nil
}
