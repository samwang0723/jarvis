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

	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/app/dto"
)

func (h *handlerImpl) GetStakeConcentration(
	ctx context.Context,
	req *dto.GetStakeConcentrationRequest,
) (*domain.StakeConcentration, error) {
	res, err := h.dataService.WithUserID(ctx).GetStakeConcentration(ctx, req)
	if err != nil {
		return nil, err
	}
	stakeConcentration := &domain.StakeConcentration{
		StockID:         res.StockID,
		Date:            res.Date,
		SumBuyShares:    res.SumBuyShares,
		SumSellShares:   res.SumSellShares,
		AvgBuyPrice:     res.AvgBuyPrice,
		AvgSellPrice:    res.AvgSellPrice,
		Concentration1:  res.Concentration1,
		Concentration5:  res.Concentration5,
		Concentration10: res.Concentration10,
		Concentration20: res.Concentration20,
		Concentration60: res.Concentration60,
	}
	stakeConcentration.ID = res.ID
	stakeConcentration.CreatedAt = res.CreatedAt
	stakeConcentration.UpdatedAt = res.UpdatedAt
	stakeConcentration.DeletedAt = res.DeletedAt

	return stakeConcentration, nil
}
