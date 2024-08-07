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

func (h *handlerImpl) ListStock(
	ctx context.Context,
	req *dto.ListStockRequest,
) (*dto.ListStockResponse, error) {
	entries, totalCount, err := h.dataService.WithUserID(ctx).ListStock(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.ListStockResponse{
		Offset:     req.Offset,
		Limit:      req.Limit,
		Entries:    entries,
		TotalCount: totalCount,
	}, nil
}

func (h *handlerImpl) ListCategories(ctx context.Context) (*dto.ListCategoriesResponse, error) {
	entries, err := h.dataService.WithUserID(ctx).ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]string, len(entries))
	for i, v := range entries {
		resp[i] = v
	}

	return &dto.ListCategoriesResponse{
		Entries: resp,
	}, nil
}
