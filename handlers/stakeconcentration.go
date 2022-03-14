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

package handlers

import (
	"context"
	"fmt"

	"github.com/samwang0723/jarvis/dto"
	"github.com/samwang0723/jarvis/entity"
)

func (h *handlerImpl) CreateStakeConcentration(ctx context.Context, req *dto.CreateStakeConcentrationRequest) (*dto.CreateStakeConcentrationResponse, error) {
	err := h.dataService.CreateStakeConcentration(ctx, req)
	if err != nil {
		return nil, err
	}
	res, err := h.dataService.GetStakeConcentration(ctx, &dto.GetStakeConcentrationRequest{StockID: req.StockID})
	if err != nil {
		return nil, err
	}
	return &dto.CreateStakeConcentrationResponse{Entry: res}, nil
}

func (h *handlerImpl) RefreshStakeConcentration(ctx context.Context, stockId string, date string) error {
	if !h.dataService.HasStakeConcentration(ctx, date) {
		return fmt.Errorf("not valid date for concentration")
	}
	concentrations := &[]interface{}{}
	c := h.calculateConcentration(ctx, stockId, date)
	if c == nil {
		return fmt.Errorf("failed to calculate concentration of %s, %s", stockId, date)
	}
	*concentrations = append(*concentrations, c)
	return h.dataService.BatchUpdateStakeConcentration(ctx, concentrations)
}

func (h *handlerImpl) calculateConcentration(ctx context.Context, stockId string, date string) *entity.StakeConcentration {
	m, err := h.dataService.GetStakeConcentrationsWithVolumes(ctx, stockId, date)
	if err != nil || len(m) < 5 {
		return nil
	}
	return &entity.StakeConcentration{
		StockID:          stockId,
		Date:             date,
		Concentration_1:  m[1],
		Concentration_5:  m[5],
		Concentration_10: m[10],
		Concentration_20: m[20],
		Concentration_60: m[60],
	}
}

func (h *handlerImpl) GetStakeConcentration(ctx context.Context, req *dto.GetStakeConcentrationRequest) (*entity.StakeConcentration, error) {
	res, err := h.dataService.GetStakeConcentration(ctx, req)
	if err != nil {
		return nil, err
	}
	stakeConcentration := &entity.StakeConcentration{
		StockID:          res.StockID,
		Date:             res.Date,
		SumBuyShares:     res.SumBuyShares,
		SumSellShares:    res.SumSellShares,
		AvgBuyPrice:      res.AvgBuyPrice,
		AvgSellPrice:     res.AvgSellPrice,
		Concentration_1:  res.Concentration_1,
		Concentration_5:  res.Concentration_5,
		Concentration_10: res.Concentration_10,
		Concentration_20: res.Concentration_20,
		Concentration_60: res.Concentration_60,
	}
	stakeConcentration.ID = res.ID
	stakeConcentration.CreatedAt = res.CreatedAt
	stakeConcentration.UpdatedAt = res.UpdatedAt
	stakeConcentration.DeletedAt = res.DeletedAt

	return stakeConcentration, nil
}
