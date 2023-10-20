package handlers

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (h *handlerImpl) CreateTransaction(
	ctx context.Context,
	req *dto.CreateTransactionRequest,
) (*dto.CreateTransactionResponse, error) {
	debitAmount, creditAmount := float32(0.0), float32(0.0)
	switch req.OrderType {
	case entity.OrderTypeDeposit:
		creditAmount = req.Amount
	case entity.OrderTypeWithdraw:
		debitAmount = req.Amount
	default:
		return &dto.CreateTransactionResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: "invalid order type",
			Success:      false,
		}, errOrderTypeNotAllowed
	}

	err := h.dataService.CreateTransaction(ctx, req.OrderType, creditAmount, debitAmount)
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
