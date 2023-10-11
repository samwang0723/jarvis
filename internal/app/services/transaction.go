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
	"context"
	"errors"

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

var ErrNoTransactionToCreate = errors.New("no transaction to create")

func (s *serviceImpl) CreateTransactions(ctx context.Context, obj *entity.Transaction) error {
	transactions := s.chainTransactions(obj, obj.OriginalExchangeDate)
	if len(transactions) == 0 {
		return ErrNoTransactionToCreate
	}

	err := s.dal.CreateTransactions(ctx, transactions)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create transaction")

		return err
	}

	return nil
}

func (s *serviceImpl) GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	transaction, err := s.dal.GetTransactionByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get transaction by id")

		return nil, err
	}

	return transaction, nil
}

func (s *serviceImpl) ListTransactions(
	ctx context.Context,
	req *dto.ListTransactionsRequest,
) (objs []*entity.Transaction, totalCount int64, err error) {
	transactions, totalCount, err := s.dal.ListTransactions(ctx, req.UserID, req.StartDate, req.EndDate)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list transactions")

		return nil, 0, err
	}

	return transactions, totalCount, nil
}

func (s *serviceImpl) chainTransactions(
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

	tax := s.taxCalculation(source, originalExchangeDate)
	if tax != nil {
		output = append(output, tax)
	}

	fee := s.feeCalculation(source)
	if fee != nil {
		output = append(output, fee)
	}

	return output
}

func (s *serviceImpl) taxCalculation(
	source *entity.Transaction,
	originalExchangeDate string,
) (output *entity.Transaction) {
	if source.ReferenceID != nil {
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

func (s *serviceImpl) feeCalculation(source *entity.Transaction) (output *entity.Transaction) {
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
