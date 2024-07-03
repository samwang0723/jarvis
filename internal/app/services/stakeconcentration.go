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
	"errors"
	"math"

	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/entity"
)

const (
	volumeCount  = 1000
	twoDigits    = 100
	roundDecimal = 10
	index1       = 0
	index5       = 4
	index10      = 9
	index20      = 19
	index60      = 59
)

var errCannotCastStakeConcentration = errors.New("cannot cast interface to *dto.StakeConcentration")

func (s *serviceImpl) GetStakeConcentration(
	ctx context.Context,
	req *dto.GetStakeConcentrationRequest,
) (*entity.StakeConcentration, error) {
	return s.dal.GetStakeConcentrationByStockID(ctx, req.StockID, req.Date)
}

func (s *serviceImpl) BatchUpsertStakeConcentration(ctx context.Context, objs *[]any) error {
	// Replicate the value from interface to *entity.StakeConcentration
	stakeConcentrations := []*entity.StakeConcentration{}
	for _, v := range *objs {
		if val, ok := v.(*entity.StakeConcentration); ok {
			s.calculateConcentration(ctx, val)
			stakeConcentrations = append(stakeConcentrations, val)
		} else {
			return errCannotCastStakeConcentration
		}
	}

	return s.dal.BatchUpsertStakeConcentration(ctx, stakeConcentrations)
}

//nolint:nolintlint, cyclop
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
			op := float32(0.0)
			if ref.Diff[cursor] != 0 && sumTradeShares > 0 {
				p := (float64(ref.Diff[cursor]) / float64(sumTradeShares/volumeCount)) * twoDigits
				op = float32(math.Round(p*roundDecimal) / roundDecimal)
			}

			switch idx {
			case index1:
				ref.Concentration1 = op
			case index5:
				ref.Concentration5 = op
			case index10:
				ref.Concentration10 = op
			case index20:
				ref.Concentration20 = op
			case index60:
				ref.Concentration60 = op
			}
			cursor++
		}
	}
}
