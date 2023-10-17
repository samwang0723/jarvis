package handlers

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (h *handlerImpl) CreateOrder(
	ctx context.Context,
	req *dto.CreateOrderRequest,
) (*dto.CreateOrderResponse, error) {
	if req.OrderType != entity.OrderTypeBuy && req.OrderType != entity.OrderTypeSell {
		h.logger.Error().Err(errOrderTypeNotAllowed).Msg("invalid order type")

		return &dto.CreateOrderResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: "invalid order type",
			Success:      false,
		}, errOrderTypeNotAllowed
	}

	order, err := entity.NewOrder(
		req.UserID,
		req.OrderType,
		req.StockID,
		req.ExchangeDate,
		req.TradePrice,
		req.Quantity,
	)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create order")

		return &dto.CreateOrderResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}, err
	}

	err = h.dataService.CreateOrder(ctx, order, req.OrderType)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create order")

		return &dto.CreateOrderResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: err.Error(),
			Success:      false,
		}, err
	}

	return &dto.CreateOrderResponse{
		Status:       dto.StatusSuccess,
		ErrorCode:    "",
		ErrorMessage: "",
		Success:      true,
	}, nil
}
