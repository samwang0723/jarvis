package services

import (
	"context"
	"fmt"

	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/helper"
)

const (
	closeToHighestToday = 0.985
	realtimeVolume      = 3000
)

// Define the SelectionStrategy interface
type SelectionStrategy interface {
	ListSelections(ctx context.Context, req *dto.ListSelectionRequest) ([]*domain.Selection, error)
}

// RealTimeSelectionStrategy implements the SelectionStrategy interface for real-time selections
type RealTimeSelectionStrategy struct {
	service *serviceImpl
}

func (r *RealTimeSelectionStrategy) ListSelections(
	ctx context.Context,
	req *dto.ListSelectionRequest,
) ([]*domain.Selection, error) {
	objs, err := r.service.ListRealTimeSelections(ctx, req)
	if err != nil {
		r.service.logger.Error().Err(err).Msg("list realtime selections")
		return nil, err
	}
	return objs, nil
}

// HistoricalSelectionStrategy implements the SelectionStrategy interface for historical selections
type HistoricalSelectionStrategy struct {
	service *serviceImpl
}

func (h *HistoricalSelectionStrategy) ListSelections(
	ctx context.Context,
	req *dto.ListSelectionRequest,
) ([]*domain.Selection, error) {
	selections, err := h.service.dal.ListSelections(ctx, req.Date, req.Strict)
	if err != nil {
		h.service.logger.Error().Err(err).Msg("list selections data record retrieval")
		return nil, err
	}

	objs, err := h.service.executeAnalysisEngine(ctx, selections, req.Strict, req.Date)
	if err != nil {
		h.service.logger.Error().Err(err).Msg("list selections advanced filtering")
		return nil, err
	}

	return objs, nil
}

func (s *serviceImpl) ListSelections(
	ctx context.Context,
	req *dto.ListSelectionRequest,
) ([]*domain.Selection, error) {
	var strategy SelectionStrategy

	today := helper.Today()
	if req.Date == today {
		// check newest concentration date
		latestDate := s.dal.GetStakeConcentrationLatestDataPoint(ctx)
		if req.Date != latestDate {
			// Case 1: req.Date is today and no data in DAL layer
			strategy = &RealTimeSelectionStrategy{service: s}
		} else {
			// Case 2: req.Date is today and has data in DAL layer
			strategy = &HistoricalSelectionStrategy{service: s}
		}
	} else if req.Date < today {
		// Case 3: req.Date is not today and must be less than today
		strategy = &HistoricalSelectionStrategy{service: s}
	} else {
		return nil, fmt.Errorf("req.Date cannot be in the future")
	}

	return strategy.ListSelections(ctx, req)
}

func (s *serviceImpl) latestStockStatSnapshot(
	ctx context.Context,
) (map[string]*domain.Selection, error) {
	m := make(map[string]*domain.Selection)
	snapshot, err := s.dal.LatestStockStatSnapshot(ctx)
	if err != nil {
		return nil, err
	}

	for _, c := range snapshot {
		m[c.StockID] = c
	}

	return m, nil
}
