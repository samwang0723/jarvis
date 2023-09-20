// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dal

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
)

// CreateBalanceView creates a balance view while user record is created.
func (i *dalImpl) CreateBalanceView(ctx context.Context, obj *entity.BalanceView) error {
	err := i.db.Create(obj).Error

	return err
}

// UpdateBalanceView updates a balance view while user inserts a transaction.
func (i *dalImpl) UpdateBalanceView(ctx context.Context, obj *entity.BalanceView) error {
	err := i.db.Unscoped().Model(&entity.BalanceView{}).Save(obj).Error

	return err
}

func (i *dalImpl) GetBalanceViewByUserID(ctx context.Context, userID entity.ID) (*entity.BalanceView, error) {
	res := &entity.BalanceView{}
	if err := i.db.First(res, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return res, nil
}
