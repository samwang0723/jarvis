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
package handlers

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/dto"
)

func (h *handlerImpl) ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) (*dto.ListDailyCloseResponse, error) {
	entries, totalCount, err := h.dataService.ListDailyClose(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, obj := range entries {
		avg, err := h.dataService.GetAverages(ctx, obj.StockID, obj.Date)
		if err != nil {
			return nil, err
		}
		obj.Average = avg
	}

	return &dto.ListDailyCloseResponse{
		Offset:     req.Offset,
		Limit:      req.Limit,
		Entries:    entries,
		TotalCount: totalCount,
	}, nil
}
