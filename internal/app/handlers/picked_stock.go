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
	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (h *handlerImpl) InsertPickedStocks(
	ctx context.Context,
	req *dto.InsertPickedStocksRequest,
) (*dto.InsertPickedStocksResponse, error) {
	objs := []*entity.PickedStock{}
	for _, stock := range req.StockIDs {
		objs = append(objs, &entity.PickedStock{
			StockID: stock,
		})
	}

	err := h.dataService.BatchUpsertPickedStocks(ctx, objs)
	if err != nil {
		return &dto.InsertPickedStocksResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}, err
	}

	return &dto.InsertPickedStocksResponse{
		Status:       dto.StatusSuccess,
		ErrorCode:    "",
		ErrorMessage: "",
		Success:      true,
	}, nil
}

func (h *handlerImpl) DeletePickedStocks(
	ctx context.Context,
	req *dto.DeletePickedStocksRequest,
) (*dto.DeletePickedStocksResponse, error) {
	err := h.dataService.DeletePickedStockByID(ctx, req.StockID)
	if err != nil {
		return &dto.DeletePickedStocksResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}, err
	}

	return &dto.DeletePickedStocksResponse{
		Status:       dto.StatusSuccess,
		ErrorCode:    "",
		ErrorMessage: "",
		Success:      true,
	}, nil
}

func (h *handlerImpl) ListPickedStocks(ctx context.Context) (*dto.ListPickedStocksResponse, error) {
	entries, err := h.dataService.ListPickedStock(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.ListPickedStocksResponse{
		Entries: entries,
	}, nil
}
