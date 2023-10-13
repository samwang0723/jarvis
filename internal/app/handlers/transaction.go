package handlers

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
)

const (
	taiwanStockQuantity = 1000
)

func (h *handlerImpl) CreateTransaction(
	ctx context.Context,
	req *dto.CreateTransactionRequest,
) (*dto.CreateTransactionResponse, error) {
	debitAmount, creditAmount := float32(0.0), float32(0.0)
	switch req.OrderType {
	case entity.OrderTypeBid:
		debitAmount = req.TradePrice * float32(req.Quantity) * taiwanStockQuantity
	case entity.OrderTypeAsk:
		creditAmount = req.TradePrice * float32(req.Quantity) * taiwanStockQuantity
	}

	transaction, err := entity.NewTransaction(
		req.StockID,
		req.UserID,
		req.OrderType,
		req.TradePrice,
		req.Quantity,
		req.ExchangeDate,
		creditAmount,
		debitAmount,
		req.Description,
		req.ReferenceID,
	)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create transaction")

		return &dto.CreateTransactionResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}, err
	}

	if req.OriginalExchangeDate != "" {
		transaction.OriginalExchangeDate = req.OriginalExchangeDate
	}

	err = h.dataService.CreateTransaction(ctx, transaction)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create transaction")

		return &dto.CreateTransactionResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}, err
	}

	return &dto.CreateTransactionResponse{
		Status:       dto.StatusSuccess,
		ErrorCode:    "",
		ErrorMessage: "",
		Success:      true,
	}, nil
}
