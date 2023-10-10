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
package entity

const (
	EventTransactionPending    = "TransactionPending"
	EventTransactionProcessing = "TransactionProcessing"
	EventTransactionCompleted  = "TransactionCompleted"
	EventTransactionFailed     = "TransactionFailed"
	EventTransactionCancelled  = "TransactionCancelled"
)

type Event struct {
	EventSourcingModel

	AggregateID uint64 `gorm:"column:aggregate_id"` // foreign key to the Transaction table
	EventType   string `gorm:"column:event_type"`   // e.g., "TransactionCreated", "TransactionStatusUpdated"
	Payload     string `gorm:"column:payload"`      // JSON-encoded event data
	Version     int    `gorm:"column:version"`      // event version number, used for ordering events
}
