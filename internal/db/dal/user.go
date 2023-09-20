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

func (i *dalImpl) CreateUser(ctx context.Context, obj *entity.User) error {
	err := i.db.Create(obj).Error

	return err
}

func (i *dalImpl) UpdateUser(ctx context.Context, obj *entity.User) error {
	err := i.db.Unscoped().Model(&entity.User{}).Save(obj).Error

	return err
}

func (i *dalImpl) DeleteUserByID(ctx context.Context, id entity.ID) error {
	err := i.db.Delete(&entity.User{}, id).Error

	return err
}

func (i *dalImpl) GetUserByID(ctx context.Context, id uint64) (*entity.User, error) {
	res := &entity.User{}
	if err := i.db.First(res, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (i *dalImpl) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	res := &entity.User{}
	if err := i.db.First(res, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (i *dalImpl) GetUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	res := &entity.User{}
	if err := i.db.First(res, "phone = ?", phone).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (i *dalImpl) ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	res := []*entity.User{}
	if err := i.db.Offset(offset).Limit(limit).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}
