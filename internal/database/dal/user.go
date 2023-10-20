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
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/database"
	"gorm.io/gorm"
)

const (
	sessionExpiredDays = 5
)

var errNoUserID = errors.New("no user id used in update method")

func (i *dalImpl) CreateUser(ctx context.Context, obj *entity.User) error {
	err := i.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(obj).Error; err != nil {
			return err
		}

		ctx = database.WithTx(ctx, tx)
		balance, err := entity.NewBalanceView(obj.ID.Uint64(), 0.0)
		if err != nil {
			return err
		}

		if err := i.balanceRepository.Save(ctx, balance); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (i *dalImpl) UpdateUser(ctx context.Context, obj *entity.User) error {
	if obj.ID.Uint64() == 0 {
		return errNoUserID
	}

	err := i.db.Unscoped().Omit("SessionID", "SessionExpiredAt").Model(&entity.User{}).Save(obj).Error

	return err
}

func (i *dalImpl) UpdateSessionID(ctx context.Context, obj *entity.User) error {
	if obj.ID.Uint64() == 0 {
		return errNoUserID
	}

	var err error
	if obj.SessionExpiredAt == nil || time.Now().After(*obj.SessionExpiredAt) {
		obj.SessionID = uuid.New().String()
		err = i.db.Raw(`
                        UPDATE users 
                        SET session_id = ?, session_expired_at = ? 
                        WHERE id = ?
                `, obj.SessionID, time.Now().AddDate(0, 0, sessionExpiredDays), obj.ID.Uint64()).Error
	}

	return err
}

func (i *dalImpl) DeleteUserByID(ctx context.Context, id uint64) error {
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

func (i *dalImpl) ListUsers(
	ctx context.Context,
	offset,
	limit int32,
) (objs []*entity.User, totalCount int64, err error) {
	sql := "select count(*) from users where deleted_at is null"
	err = i.db.Raw(sql).Scan(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	sql = fmt.Sprintf(`select * from users where deleted_at is null order by created_at desc limit %d, %d`, offset, limit)
	if err := i.db.Raw(sql).Scan(&objs).Error; err != nil {
		return nil, totalCount, err
	}

	return objs, totalCount, nil
}
