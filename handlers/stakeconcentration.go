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
	"math"

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

func (h *handlerImpl) RefreshStakeConcentration(ctx context.Context, req *dto.RefreshStakeConcentrationRequest) (*dto.RefreshStakeConcentrationResponse, error) {
	stockId := req.StockID
	date := req.Date
	if !h.dataService.HasStakeConcentration(ctx, date) {
		return &dto.RefreshStakeConcentrationResponse{
			Code:     400,
			Error:    "Bad Request",
			Messages: fmt.Sprintf("No valid date for concentration: %s, %s", stockId, date),
		}, fmt.Errorf("No valid date for concentration: %s, %s", stockId, date)
	}
	concentrations := &[]interface{}{}
	c := h.calculateConcentration(ctx, stockId, date, req.Diff)
	if c == nil {
		return &dto.RefreshStakeConcentrationResponse{
			Code:     400,
			Error:    "Bad Request",
			Messages: fmt.Sprintf("Failed to calculate concentration: %s, %s", stockId, date),
		}, fmt.Errorf("Failed to calculate concentration: %s, %s", stockId, date)
	}
	*concentrations = append(*concentrations, c)
	err := h.dataService.BatchUpdateStakeConcentration(ctx, concentrations)
	if err != nil {
		return &dto.RefreshStakeConcentrationResponse{
			Code:     500,
			Error:    "Internal Server Error",
			Messages: fmt.Sprintf("Failed to upate concentration (%s): %s, %s", err, stockId, date),
		}, fmt.Errorf("Failed to update concentration(%s): %s, %s", err, stockId, date)
	}

	return &dto.RefreshStakeConcentrationResponse{
		Code:     200,
		Messages: fmt.Sprintf("Concentration updated successfully: %s, %s", stockId, date),
	}, nil
}

func (h *handlerImpl) calculateConcentration(ctx context.Context, stockId string, date string, diff []int32) *entity.StakeConcentration {
	bases, err := h.dataService.GetStakeConcentrationsWithVolumes(ctx, stockId, date)
	if err != nil || len(bases) < 60 {
		return nil
	}

	resp := &entity.StakeConcentration{
		StockID: stockId,
		Date:    date,
	}

	sumTradeShares := uint64(0)
	cursor := 0
	for idx, c := range bases {
		sumTradeShares += c.TradeShares

		if idx == 0 || idx == 4 || idx == 9 || idx == 19 || idx == 59 {
			p, op := 0.0, float32(0.0)
			if diff[cursor] != 0 && sumTradeShares > 0 {
				p = (float64(diff[cursor]) / float64(sumTradeShares/1000)) * 100
				op = float32(math.Round(p*10) / 10)
			}

			switch idx {
			case 0:
				resp.Concentration_1 = op
			case 4:
				resp.Concentration_5 = op
			case 9:
				resp.Concentration_10 = op
			case 19:
				resp.Concentration_20 = op
			case 59:
				resp.Concentration_60 = op
			}
			cursor++
		}
	}
	return resp
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
