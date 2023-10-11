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

import jsoniter "github.com/json-iterator/go"

const (
	OrderTypeBid = iota
	OrderTypeAsk
	OrderTypeFee
	OrderTypeTax
	OrderTypeLending
)

type TransactionPayload struct {
	CreditAmount float32 `json:"creditAmount"`
	DebitAmount  float32 `json:"debitAmount"`
	EventType    string  `json:"type"`
	Auditor      string  `json:"auditor"`
	Description  string  `json:"description"`
}

func (payload TransactionPayload) ToJSON() string {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	jsonString, err := json.MarshalToString(payload)
	if err != nil {
		return ""
	}

	return jsonString
}

type Transaction struct {
	EventSourcingModel

	StockID              string  `gorm:"column:stock_id" json:"stockId"`
	UserID               uint64  `gorm:"column:user_id" json:"userId"`
	OrderType            int32   `gorm:"column:order_type" json:"orderType"`
	TradePrice           float32 `gorm:"column:trade_price" json:"tradePrice"`
	Quantity             uint64  `gorm:"column:quantity" json:"quantity"`
	ExchangeDate         string  `gorm:"column:exchange_date" json:"exchangeDate"`
	CreditAmount         float32 `gorm:"column:credit_amount" json:"creditAmount"`
	DebitAmount          float32 `gorm:"column:debit_amount" json:"debitAmount"`
	Description          string  `gorm:"column:description" json:"description"`
	ReferenceID          *uint64 `gorm:"column:reference_id" json:"referenceId"`
	Status               string  `gorm:"column:status" json:"status,omitempty"`
	OriginalExchangeDate string
}

func (Transaction) TableName() string {
	return "transactions"
}
