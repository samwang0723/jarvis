package handlers

import (
	"context"
	"samwang0723/jarvis/dto"
)

func (h *handlerImpl) ListDailyClose(ctx context.Context, req *dto.ListDailyCloseRequest) (*dto.ListDailyCloseResponse, error) {
	entries, totalCount, err := h.dataService.ListDailyClose(ctx, req)
	if err != nil {
		return nil, err
	}
	return &dto.ListDailyCloseResponse{
		Offset:     req.Offset,
		Limit:      req.Limit,
		Entries:    entries,
		TotalCount: int(totalCount),
	}, nil
}
