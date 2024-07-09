package handlers

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
)

func (h *handlerImpl) CreateOrder(
	ctx context.Context,
	req *dto.CreateOrderRequest,
) (*dto.CreateOrderResponse, error) {
	if req.OrderType != domain.OrderTypeBuy && req.OrderType != domain.OrderTypeSell {
		h.logger.Error().Err(errOrderTypeNotAllowed).Msg("invalid order type")

		return &dto.CreateOrderResponse{
			Status:       dto.StatusError,
			ErrorCode:    "",
			ErrorMessage: "invalid order type",
			Success:      false,
		}, errOrderTypeNotAllowed
	}

	err := h.dataService.WithUserID(ctx).CreateOrder(ctx, req)
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

func (h *handlerImpl) ListOrders(
	ctx context.Context,
	req *dto.ListOrderRequest,
) (*dto.ListOrderResponse, error) {
	orders, totalCount, err := h.dataService.WithUserID(ctx).ListOrders(ctx, req)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list orders")

		return nil, err
	}

	return &dto.ListOrderResponse{
		Entries:    orders,
		Offset:     req.Offset,
		Limit:      req.Limit,
		TotalCount: totalCount,
	}, nil
}
