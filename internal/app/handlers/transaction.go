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

func (h *handlerImpl) ListTransactions(
	ctx context.Context,
	req *dto.ListTransactionsRequest,
) (*dto.ListTransactionsResponse, error) {
	entries, totalCount, err := h.dataService.ListTransactions(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.ListTransactionsResponse{
		Entries:    entries,
		TotalCount: totalCount,
	}, nil
}

func (h *handlerImpl) CreateTransactions(
	ctx context.Context,
	req *dto.CreateTransactionsRequest,
) (*dto.CreateTransactionsResponse, error) {
	if req.OrderType == entity.OrderTypeAsk && req.ReferenceID == 0 {
		return &dto.CreateTransactionsResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: "system not support short selling now",
			Success:      false,
		}, nil
	}

	transaction := &entity.Transaction{
		StockID:      req.StockID,
		UserID:       req.UserID,
		OrderType:    req.OrderType,
		TradePrice:   req.TradePrice,
		Quantity:     req.Quantity,
		ExchangeDate: req.ExchangeDate,
		Description:  req.Description,
	}

	if req.ReferenceID == 0 {
		transaction.ReferenceID = nil
	} else {
		transaction.ReferenceID = &req.ReferenceID
	}

	if req.OriginalExchangeDate != "" {
		transaction.OriginalExchangeDate = req.OriginalExchangeDate
	}

	err := h.dataService.CreateTransactions(ctx, transaction)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create transaction")

		return &dto.CreateTransactionsResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}, err
	}

	return &dto.CreateTransactionsResponse{
		Status:       dto.StatusSuccess,
		ErrorCode:    "",
		ErrorMessage: "",
		Success:      true,
	}, nil
}
