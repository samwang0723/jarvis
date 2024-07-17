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
package handlers

import (
	"context"
	"time"

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/cache"
)

const (
	cronLockPeriod = 1
)

func (h *handlerImpl) ListSelections(
	ctx context.Context,
	req *dto.ListSelectionRequest,
) (*dto.ListSelectionResponse, error) {
	entries, err := h.dataService.WithUserID(ctx).ListSelections(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.ListSelectionResponse{
		Entries: entries,
	}, nil
}

func (h *handlerImpl) CronjobPresetRealtimeMonitoringKeys(
	ctx context.Context,
	schedule string,
) error {
	err := h.dataService.AddJob(ctx, schedule, func() {
		err := h.dataService.CronjobPresetRealtimeMonitoringKeys(ctx)
		if err != nil {
			h.logger.Error().Msgf("failed to preset real time keys: %s", err)
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *handlerImpl) CrawlingRealTimePrice(ctx context.Context, schedule string) error {
	err := h.dataService.AddJob(ctx, schedule, func() {
		if h.dataService.ObtainLock(ctx, cache.CronjobLock, cronLockPeriod*time.Minute) == nil {
			h.logger.Info().Msg("cronjob lock is not obtained")
			return
		}
		err := h.dataService.CrawlingRealTimePrice(ctx)
		if err != nil {
			h.logger.Error().Msgf("failed to retrieve real time price: %s", err)
		}
	})
	if err != nil {
		return err
	}

	return nil
}
