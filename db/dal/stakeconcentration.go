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

package dal

import (
	"context"
	"samwang0723/jarvis/entity"
)

func (i *dalImpl) CreateStakeConcentration(ctx context.Context, obj *entity.StakeConcentration) error {
	err := i.db.Create(obj).Error
	return err
}

func (i *dalImpl) GetStakeConcentrationByStockID(ctx context.Context, stockID string) (*entity.StakeConcentration, error) {
	res := &entity.StakeConcentration{}
	if err := i.db.First(res, "stock_id = ?", stockID).Error; err != nil {
		return nil, err
	}
	return res, nil
}
