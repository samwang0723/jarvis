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
	"fmt"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"gorm.io/gorm"
)

func (i *dalImpl) CreateTransactions(ctx context.Context, objs []*entity.Transaction) error {
	var totalDebitAmount float32
	var totalCreditAmount float32
	var createdReferenceID uint64

	balanceView := &entity.BalanceView{}
	if err := i.db.First(balanceView, "user_id = ?", objs[0].UserID).Error; err != nil {
		return err
	}

	err := i.db.Transaction(func(tx *gorm.DB) error {
		totalDebitAmount = 0.0
		totalCreditAmount = 0.0
		for _, obj := range objs {
			if createdReferenceID != 0 {
				obj.ReferenceID = createdReferenceID
			}
			if err := tx.Create(obj).Error; err != nil {
				return err
			}
			if createdReferenceID == 0 {
				createdReferenceID = uint64(obj.ID)
			}
			totalDebitAmount += obj.DebitAmount
			totalCreditAmount += obj.CreditAmount
		}

		balanceView.CurrentAmount = balanceView.CurrentAmount - totalDebitAmount + totalCreditAmount
		err := tx.Model(&entity.BalanceView{}).
			Where("user_id = ?", balanceView.UserID).
			Updates(balanceView).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (i *dalImpl) GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	res := &entity.Transaction{}
	if err := i.db.First(res, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (i *dalImpl) ListTransactions(
	ctx context.Context,
	userID uint64,
	limit, offset int32,
) (objs []*entity.Transaction, totalCount int64, err error) {
	sql := fmt.Sprintf(`select count(*) from transactions where user_id = %d and deleted_at is null`, userID)
	err = i.db.Raw(sql).Scan(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	sql = fmt.Sprintf(`select * from transactions where user_id = %d 
                and deleted_at is null order by created_at desc limit %d, %d`, userID, offset, limit)
	if err := i.db.Raw(sql).Scan(&objs).Error; err != nil {
		return nil, 0, err
	}

	return objs, totalCount, nil
}
