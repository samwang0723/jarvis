package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (s *serviceImpl) CreateOrder(ctx context.Context, source *entity.Order, orderType string) error {
	transactions := s.chainTransactions(source, orderType)

	return s.dal.CreateChainTransactions(ctx, transactions)
}

func (s *serviceImpl) chainTransactions(
	req *entity.Order,
	orderType string,
) (chainedTransactions []*entity.Transaction) {
	debitAmount, creditAmount := float32(0.0), float32(0.0)
	switch orderType {
	case entity.OrderTypeBuy:
		debitAmount = req.BuyPrice * float32(req.BuyQuantity) * taiwanStockQuantity
	case entity.OrderTypeSell:
		creditAmount = req.SellPrice * float32(req.SellQuantity) * taiwanStockQuantity
	}

	transaction, err := entity.NewTransaction(
		req.UserID,
		orderType,
		creditAmount,
		debitAmount,
		req.ID,
	)
	if err != nil {
		return nil
	}

	chainedTransactions = append(chainedTransactions, transaction)

	tax := s.genTaxTransaction(req, orderType)
	if tax != nil {
		chainedTransactions = append(chainedTransactions, tax)
	}

	fee := s.genFeeTransaction(req, orderType)
	if fee != nil {
		chainedTransactions = append(chainedTransactions, fee)
	}

	return chainedTransactions
}

func (s *serviceImpl) genTaxTransaction(req *entity.Order, orderType string) *entity.Transaction {
	// only charge tax on partial order close or complete order close
	if req.BuyPrice > 0 && req.SellPrice > 0 {
		debitAmount := float32(0.0)
		if orderType == entity.OrderTypeBuy {
			debitAmount = req.BuyPrice * float32(req.BuyQuantity) * taiwanStockQuantity * taxRate
		} else if orderType == entity.OrderTypeSell {
			debitAmount = req.SellPrice * float32(req.SellQuantity) * taiwanStockQuantity * taxRate
		}
		if req.BuyExchangeDate == req.SellExchangeDate {
			debitAmount *= dayTradeTaxRate
		}

		//nolint: errcheck // return nil transaction
		output, _ := entity.NewTransaction(
			req.UserID,
			entity.OrderTypeTax,
			0,
			debitAmount,
			req.ID,
		)

		return output
	}

	return nil
}

func (s *serviceImpl) genFeeTransaction(req *entity.Order, orderType string) *entity.Transaction {
	debitAmount := float32(0.0)
	if orderType == entity.OrderTypeBuy {
		debitAmount = req.BuyPrice * float32(req.BuyQuantity) * taiwanStockQuantity * feeRate * brokerFeeDiscount
	} else if orderType == entity.OrderTypeSell {
		debitAmount = req.SellPrice * float32(req.SellQuantity) * taiwanStockQuantity * feeRate * brokerFeeDiscount
	}
	//nolint: errcheck // return nil transaction
	output, _ := entity.NewTransaction(
		req.UserID,
		entity.OrderTypeFee,
		0,
		debitAmount,
		req.ID,
	)

	return output
}
