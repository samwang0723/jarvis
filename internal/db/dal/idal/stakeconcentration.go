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

package idal

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

type IStakeConcentrationDAL interface {
	CreateStakeConcentration(ctx context.Context, obj *entity.StakeConcentration) error
	GetStakeConcentrationByStockID(ctx context.Context, stockID string, date string) (*entity.StakeConcentration, error)
	ListBackfillStakeConcentrationStockIDs(ctx context.Context, date string) ([]string, error)
	GetStakeConcentrationsWithVolumes(ctx context.Context, stockId string, date string) (objs []*entity.CalculationBase, err error)
	BatchUpdateStakeConcentration(ctx context.Context, objs []*entity.StakeConcentration) error
	HasStakeConcentration(ctx context.Context, date string) bool
}
