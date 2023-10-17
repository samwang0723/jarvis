package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

const (
	taiwanStockQuantity = 1000
	dayTradeTaxRate     = 0.5
	taxRate             = 0.003
	feeRate             = 0.001425
	brokerFeeDiscount   = 0.25
)

func (s *serviceImpl) CreateOrder(ctx context.Context, source *entity.Order, orderType string) error {
	transactions, err := s.chainTransactions(source, orderType)
	if err != nil {
		return errUnableToChainTransactions
	}

	return s.dal.CreateOrder(ctx, source, transactions)
}

func (s *serviceImpl) chainTransactions(
	req *entity.Order,
	orderType string,
) (chainedTransactions []*entity.Transaction, err error) {
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
		return chainedTransactions, err
	}

	chainedTransactions = append(chainedTransactions, transaction)

	tax, err := s.genTaxTransaction(req, orderType)
	if err != nil {
		return chainedTransactions, err
	} else if tax != nil {
		chainedTransactions = append(chainedTransactions, tax)
	}

	fee, err := s.genFeeTransaction(req, orderType)
	if err != nil {
		return chainedTransactions, err
	} else if fee != nil {
		chainedTransactions = append(chainedTransactions, fee)
	}

	return chainedTransactions, nil
}

func (s *serviceImpl) genTaxTransaction(req *entity.Order, orderType string) (*entity.Transaction, error) {
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

		output, err := entity.NewTransaction(
			req.UserID,
			entity.OrderTypeTax,
			0,
			debitAmount,
			req.ID,
		)

		return output, err
	}

	//nolint:nilnil // this is a special case
	return nil, nil
}

func (s *serviceImpl) genFeeTransaction(req *entity.Order, orderType string) (*entity.Transaction, error) {
	debitAmount := float32(0.0)
	if orderType == entity.OrderTypeBuy {
		debitAmount = req.BuyPrice * float32(req.BuyQuantity) * taiwanStockQuantity * feeRate * brokerFeeDiscount
	} else if orderType == entity.OrderTypeSell {
		debitAmount = req.SellPrice * float32(req.SellQuantity) * taiwanStockQuantity * feeRate * brokerFeeDiscount
	}

	output, err := entity.NewTransaction(
		req.UserID,
		entity.OrderTypeFee,
		0,
		debitAmount,
		req.ID,
	)

	return output, err
}
