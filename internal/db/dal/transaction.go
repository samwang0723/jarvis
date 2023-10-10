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

	"github.com/samwang0723/jarvis/internal/app/entity"
	"gorm.io/gorm"
)

const snapshotEventThreshold = 3

var ErrInvalidDateRange = errors.New("invalid date range")

func (i *dalImpl) CreateTransactions(ctx context.Context, objs []*entity.Transaction) error {
	transactionIDs := []uint64{}
	err := i.db.Transaction(func(tx *gorm.DB) error {
		// Get the current balance view for this user
		var balanceView entity.BalanceView
		if err := tx.First(&balanceView, "user_id = ? and deleted_at is null", objs[0].UserID).Error; err != nil {
			return err
		}
		var createdReferenceID *uint64

		for _, obj := range objs {
			if createdReferenceID != nil {
				obj.ReferenceID = createdReferenceID
			}
			if err := tx.Create(obj).Error; err != nil {
				return err
			}
			if createdReferenceID == nil {
				id := uint64(obj.ID)
				createdReferenceID = &id
			}
			transactionIDs = append(transactionIDs, obj.ID.Uint64())

			// As we are not relying on external system to update status
			// can directly loop through to final completed status.
			states := []string{
				entity.EventTransactionPending,
				entity.EventTransactionProcessing,
				entity.EventTransactionCompleted,
			}
			for idx, state := range states {
				event := &entity.Event{
					AggregateID: obj.ID.Uint64(),
					EventType:   state,
					Payload: entity.TransactionPayload{
						CreditAmount: obj.CreditAmount,
						DebitAmount:  obj.DebitAmount,
						EventType:    state,
						Auditor:      "system",
						Description:  "",
					}.ToJSON(),
					Version: idx + 1,
				}

				if err := tx.Create(event).Error; err != nil {
					return err
				}

				applyEvent(event, obj, &balanceView)
			}
		}

		if err := tx.Save(&balanceView).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// update snapshot in background
	i.updateSnapshot(transactionIDs)

	return nil
}

func (i *dalImpl) updateSnapshot(transactionIDs []uint64) {
	for _, transactionID := range transactionIDs {
		go func(transactionID uint64) {
			// Get the count of events for this transaction
			var eventCount int
			if err := i.db.Raw(
				"select count(*) from events where aggregate_id = ?", transactionID,
			).Scan(&eventCount).Error; err != nil {
				return
			}
			// If there are more than 3 events, create a new snapshot
			if eventCount >= snapshotEventThreshold {
				event := entity.Event{}
				if err := i.db.Raw(
					"select * from events where aggregate_id = ? order by version desc limit 1", transactionID,
				).Scan(&event).Error; err != nil {
					return
				}
				snapshot := entity.Snapshot{
					AggregateID: transactionID,
					Data:        event.Payload,
					Version:     event.Version,
				}

				// Save the snapshot to the database
				i.db.Create(&snapshot)
			}
		}(transactionID)
	}
}

func applyEvent(event *entity.Event, transaction *entity.Transaction, balanceView *entity.BalanceView) {
	switch event.EventType {
	case entity.EventTransactionPending:
		// Do nothing
	case entity.EventTransactionProcessing:
	case entity.EventTransactionCompleted:
		balanceView.CurrentAmount += transaction.CreditAmount - transaction.DebitAmount
	case entity.EventTransactionFailed:
	case entity.EventTransactionCancelled:
	}
}

func (i *dalImpl) GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	res := &entity.Transaction{}
	if err := i.db.Raw(`
                SELECT 
                    t.*,
                    COALESCE(
                        JSON_UNQUOTE(JSON_EXTRACT(s.data, '$.type')),
                        JSON_UNQUOTE(JSON_EXTRACT(e.payload, '$.type'))
                    ) AS status
                FROM 
                    transactions t
                LEFT JOIN 
                    snapshots s ON t.id = s.aggregate_id
                LEFT JOIN (
                    SELECT 
                        aggregate_id, 
                        payload
                    FROM 
                        events
                    WHERE 
                        (aggregate_id, version) IN (
                            SELECT 
                                aggregate_id, 
                                MAX(version) 
                            FROM 
                                events 
                            GROUP BY 
                                aggregate_id
                        )
                ) e ON t.id = e.aggregate_id
                WHERE 
                    t.id = ?;
                `, id).Scan(res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (i *dalImpl) ListTransactions(
	ctx context.Context,
	userID uint64,
	startDate, endDate string,
) (objs []*entity.Transaction, totalCount int64, err error) {
	if len(startDate) < 8 || len(endDate) < 8 {
		return nil, 0, ErrInvalidDateRange
	}

	err = i.db.Raw(`select count(*) from transactions 
                where user_id = ? and created_at >= ? 
                and created_at < ?`, userID, startDate, endDate).Scan(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	if err := i.db.Raw(`
                SELECT 
                    t.*,
                    COALESCE(
                        JSON_UNQUOTE(JSON_EXTRACT(s.data, '$.type')),
                        JSON_UNQUOTE(JSON_EXTRACT(e.payload, '$.type'))
                    ) AS status
                FROM 
                    transactions t
                LEFT JOIN 
                    snapshots s ON t.id = s.aggregate_id
                LEFT JOIN (
                    SELECT 
                        aggregate_id, 
                        payload
                    FROM 
                        events
                    WHERE 
                        (aggregate_id, version) IN (
                            SELECT 
                                aggregate_id, 
                                MAX(version) 
                            FROM 
                                events 
                            GROUP BY 
                                aggregate_id
                        )
                ) e ON t.id = e.aggregate_id
                WHERE 
                    t.user_id = ?
                AND
                    t.created_at >= ?
                AND
                    t.created_at < ?
                ORDER BY 
                    t.created_at DESC;
                `, userID, startDate, endDate).Scan(&objs).Error; err != nil {
		return nil, 0, err
	}

	return objs, totalCount, nil
}
