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

const (
	TaiwanStockQuantity = 1000
	DayTradeTaxRate     = 0.5
	TaxRate             = 0.003
	FeeRate             = 0.001425
	BrokerFeeDiscount   = 0.25
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
		Offset:     req.Offset,
		Limit:      req.Limit,
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
		ReferenceID:  req.ReferenceID,
	}

	transactions := h.chainTransactions(transaction, req.OriginalExchangeDate)
	if len(transactions) == 0 {
		return &dto.CreateTransactionsResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: "no transaction to create",
			Success:      false,
		}, nil
	}

	err := h.dataService.CreateTransactions(ctx, transactions)
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

func (h *handlerImpl) chainTransactions(
	source *entity.Transaction,
	originalExchangeDate string,
) (output []*entity.Transaction) {
	switch source.OrderType {
	case entity.OrderTypeBid:
		source.DebitAmount = source.TradePrice * float32(source.Quantity) * TaiwanStockQuantity
		output = append(output, source)
	case entity.OrderTypeAsk:
		source.CreditAmount = source.TradePrice * float32(source.Quantity) * TaiwanStockQuantity
		output = append(output, source)
	}

	tax := h.taxCalculation(source, originalExchangeDate)
	if tax != nil {
		output = append(output, tax)
	}

	fee := h.feeCalculation(source)
	if fee != nil {
		output = append(output, fee)
	}

	return output
}

func (h *handlerImpl) taxCalculation(
	source *entity.Transaction,
	originalExchangeDate string,
) (output *entity.Transaction) {
	if source.ReferenceID != 0 {
		debitAmount := source.TradePrice * float32(source.Quantity) * TaiwanStockQuantity * TaxRate
		if source.ExchangeDate == originalExchangeDate {
			debitAmount *= DayTradeTaxRate
		}

		output = &entity.Transaction{
			StockID:      source.StockID,
			UserID:       source.UserID,
			OrderType:    entity.OrderTypeTax,
			TradePrice:   source.TradePrice,
			Quantity:     source.Quantity,
			ExchangeDate: source.ExchangeDate,
			Description:  source.Description,
			ReferenceID:  source.ReferenceID,
			DebitAmount:  debitAmount,
		}

		return output
	}

	return nil
}

func (h *handlerImpl) feeCalculation(source *entity.Transaction) (output *entity.Transaction) {
	debitAmount := source.TradePrice * float32(source.Quantity) * TaiwanStockQuantity * FeeRate * BrokerFeeDiscount
	output = &entity.Transaction{
		StockID:      source.StockID,
		UserID:       source.UserID,
		OrderType:    entity.OrderTypeFee,
		TradePrice:   source.TradePrice,
		Quantity:     source.Quantity,
		ExchangeDate: source.ExchangeDate,
		Description:  source.Description,
		ReferenceID:  source.ReferenceID,
		DebitAmount:  debitAmount,
	}

	return output
}
