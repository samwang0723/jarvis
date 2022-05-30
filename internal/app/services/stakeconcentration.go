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

package services

import (
	"context"
	"fmt"
	"math"
	"reflect"

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
)

func (s *serviceImpl) GetStakeConcentration(ctx context.Context, req *dto.GetStakeConcentrationRequest) (*entity.StakeConcentration, error) {
	return s.dal.GetStakeConcentrationByStockID(ctx, req.StockID, req.Date)
}

func (s *serviceImpl) BatchUpsertStakeConcentration(ctx context.Context, objs *[]interface{}) error {
	// Replicate the value from interface to *entity.StakeConcentration
	stakeConcentrations := []*entity.StakeConcentration{}
	for _, v := range *objs {
		if val, ok := v.(*entity.StakeConcentration); ok {
			s.calculateConcentration(ctx, val)
			stakeConcentrations = append(stakeConcentrations, val)
		} else {
			return fmt.Errorf("cannot cast interface to *dto.StakeConcentration: %v\n", reflect.TypeOf(v).Elem())
		}
	}

	return s.dal.BatchUpsertStakeConcentration(ctx, stakeConcentrations)
}

func (s *serviceImpl) calculateConcentration(ctx context.Context, ref *entity.StakeConcentration) {
	// pull the sum of traded volumes in order to calculate the concentration percentage
	bases, err := s.dal.GetStakeConcentrationsWithVolumes(ctx, ref.StockID, ref.Date)
	if err != nil || len(bases) < 60 {
		return
	}

	sumTradeShares := uint64(0)
	// cursor for the "diff" array, contains from 1, 5, 10, 20, 60 Buy-Sell diff records
	cursor := 0
	for idx, c := range bases {
		sumTradeShares += c.TradeShares

		if idx == 0 || idx == 4 || idx == 9 || idx == 19 || idx == 59 {
			p, op := 0.0, float32(0.0)
			if ref.Diff[cursor] != 0 && sumTradeShares > 0 {
				p = (float64(ref.Diff[cursor]) / float64(sumTradeShares/1000)) * 100
				op = float32(math.Round(p*10) / 10)
			}

			switch idx {
			case 0:
				ref.Concentration_1 = op
			case 4:
				ref.Concentration_5 = op
			case 9:
				ref.Concentration_10 = op
			case 19:
				ref.Concentration_20 = op
			case 59:
				ref.Concentration_60 = op
			}
			cursor++
		}
	}
}
