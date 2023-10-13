package services

import (
	"context"
	"errors"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

const (
	taiwanStockQuantity = 1000
	dayTradeTaxRate     = 0.5
	taxRate             = 0.003
	feeRate             = 0.001425
	brokerFeeDiscount   = 0.25
)

var errUnableToChainTransactions = errors.New("unable to create chain transactions")

func (s *serviceImpl) CreateTransaction(ctx context.Context, obj *entity.Transaction) error {
	transactions := s.chainTransactions(obj)
	if len(transactions) == 0 {
		return errUnableToChainTransactions
	}

	err := s.dal.CreateChainTransactions(ctx, transactions)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create transaction")

		return err
	}

	return nil
}

func (s *serviceImpl) chainTransactions(source *entity.Transaction) (output []*entity.Transaction) {
	output = append(output, source)

	tax := s.taxCalculation(source, source.OriginalExchangeDate)
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
) *entity.Transaction {
	if source.ReferenceID != nil {
		debitAmount := source.TradePrice * float32(source.Quantity) * taiwanStockQuantity * taxRate
		if source.ExchangeDate == originalExchangeDate {
			debitAmount *= dayTradeTaxRate
		}

		//nolint: errcheck // return nil transaction
		output, _ := entity.NewTransaction(
			source.StockID,
			source.UserID,
			entity.OrderTypeTax,
			source.TradePrice,
			source.Quantity,
			source.ExchangeDate,
			0, debitAmount,
			source.Description,
			source.ReferenceID,
		)

		return output
	}

	return nil
}

func (s *serviceImpl) feeCalculation(source *entity.Transaction) *entity.Transaction {
	debitAmount := source.TradePrice * float32(source.Quantity) * taiwanStockQuantity * feeRate * brokerFeeDiscount
	//nolint: errcheck // return nil transaction
	output, _ := entity.NewTransaction(
		source.StockID,
		source.UserID,
		entity.OrderTypeFee,
		source.TradePrice,
		source.Quantity,
		source.ExchangeDate,
		0, debitAmount,
		source.Description,
		source.ReferenceID,
	)

	return output
}
