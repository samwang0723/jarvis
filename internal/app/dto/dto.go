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
package dto

import (
	"github.com/samwang0723/jarvis/internal/app/entity"
)

const (
	StatusSuccess      = 200
	StatusBadRequest   = 400
	StatusUnauthorized = 401
	StatusError        = 500
)

type ListDailyCloseSearchParams struct {
	StockID string  `json:"stockID"`
	End     *string `json:"end,omitempty"`
	Start   string  `json:"start"`
}

type ListDailyCloseRequest struct {
	SearchParams *ListDailyCloseSearchParams `json:"searchParams"`
	Offset       int32                       `json:"offset"`
	Limit        int32                       `json:"limit"`
}

type ListDailyCloseResponse struct {
	Entries    []*entity.DailyClose `json:"entries"`
	Offset     int32                `json:"offset"`
	Limit      int32                `json:"limit"`
	TotalCount int64                `json:"totalCount"`
}

type ListThreePrimarySearchParams struct {
	End     *string `json:"end,omitempty"`
	StockID string  `json:"stockID,omitempty"`
	Start   string  `json:"start"`
}

type ListThreePrimaryRequest struct {
	SearchParams *ListThreePrimarySearchParams `json:"searchParams"`
	Offset       int32                         `json:"offset"`
	Limit        int32                         `json:"limit"`
}

type ListThreePrimaryResponse struct {
	Entries    []*entity.ThreePrimary `json:"entries"`
	Offset     int32                  `json:"offset"`
	Limit      int32                  `json:"limit"`
	TotalCount int64                  `json:"totalCount"`
}

type ListStockSearchParams struct {
	StockIDs *[]string `json:"stockIDs,omitempty"`
	Name     *string   `json:"name,omitempty"`
	Category *string   `json:"category,omitempty"`
	Country  string    `json:"country"`
}

type ListStockRequest struct {
	SearchParams *ListStockSearchParams `json:"searchParams"`
	Offset       int32                  `json:"offset"`
	Limit        int32                  `json:"limit"`
}

type ListStockResponse struct {
	Entries    []*entity.Stock `json:"entries"`
	Offset     int32           `json:"offset"`
	Limit      int32           `json:"limit"`
	TotalCount int64           `json:"totalCount"`
}

type ListCategoriesResponse struct {
	Entries []string `json:"entries"`
}

type GetStakeConcentrationRequest struct {
	StockID string `json:"stockID"`
	Date    string `json:"date"`
}

type ListSelectionRequest struct {
	Date   string `json:"date"`
	Strict bool   `json:"strict"`
}

type ListSelectionResponse struct {
	Entries []*entity.Selection `json:"entries"`
}

type ListPickedStocksResponse struct {
	Entries []*entity.Selection `json:"entries"`
}

type InsertPickedStocksRequest struct {
	StockIDs []string `json:"stockIDs"`
}

type InsertPickedStocksResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Success      bool   `json:"success"`
	Status       int    `json:"status"`
}

type DeletePickedStocksRequest struct {
	StockID string `json:"stockID"`
}

type DeletePickedStocksResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Success      bool   `json:"success"`
	Status       int    `json:"status"`
}

type ListUsersRequest struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

type ListUsersResponse struct {
	Entries    []*entity.User `json:"entries"`
	Offset     int32          `json:"offset"`
	Limit      int32          `json:"limit"`
	TotalCount int64          `json:"totalCount"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	AccessToken  string `json:"access_token"`
	Status       int    `json:"status"`
	Success      bool   `json:"success"`
}

type LogoutResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Success      bool   `json:"success"`
	Status       int    `json:"status"`
}

type CreateUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type CreateUserResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Success      bool   `json:"success"`
	Status       int    `json:"status"`
}

type GetBalanceViewRequest struct{}

type GetBalanceViewResponse struct {
	Balance *entity.BalanceView `json:"balance"`
}

type CreateTransactionRequest struct {
	OrderType string  `json:"orderType"`
	Amount    float32 `json:"amount"`
}

type CreateTransactionResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Success      bool   `json:"success"`
	Status       int    `json:"status"`
}

type CreateOrderRequest struct {
	OrderType    string  `json:"orderType"`
	StockID      string  `json:"stockID"`
	ExchangeDate string  `json:"exchangeDate"`
	TradePrice   float32 `json:"tradePrice"`
	Quantity     uint64  `json:"quantity"`
}

type CreateOrderResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Success      bool   `json:"success"`
	Status       int    `json:"status"`
}

type ListOrderSearchParams struct {
	StockIDs      *[]string `json:"stockIDs,omitempty"`
	ExchangeMonth *string   `json:"exchangeMonth,omitempty"`
	Status        *string   `json:"status,omitempty"`
}

type ListOrderRequest struct {
	SearchParams *ListOrderSearchParams `json:"searchParams"`
	Offset       int32                  `json:"offset"`
	Limit        int32                  `json:"limit"`
}

type ListOrderResponse struct {
	Entries    []*entity.Order `json:"entries"`
	Offset     int32           `json:"offset"`
	Limit      int32           `json:"limit"`
	TotalCount int64           `json:"totalCount"`
}
