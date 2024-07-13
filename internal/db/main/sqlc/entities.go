// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlcdb

import (
	"database/sql"
	"time"

	"github.com/ericlagergren/decimal"
	uuid "github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type BalanceEvent struct {
	AggregateID uuid.UUID
	ParentID    uuid.UUID
	EventType   string
	Payload     []byte
	Version     int32
	CreatedAt   time.Time
}

type BalanceView struct {
	ID        uuid.UUID
	Balance   decimal.Big
	Available decimal.Big
	Pending   decimal.Big
	Version   int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DailyClose struct {
	ID           uuid.UUID
	StockID      string
	ExchangeDate string
	TradeShares  sql.NullInt64
	Transactions sql.NullInt64
	Turnover     sql.NullInt64
	Open         decimal.Big
	Close        decimal.Big
	High         decimal.Big
	Low          decimal.Big
	PriceDiff    decimal.Big
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    sql.NullTime
}

type Order struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	StockID          string
	BuyPrice         decimal.Big
	BuyQuantity      int64
	BuyExchangeDate  string
	SellPrice        decimal.Big
	SellQuantity     int64
	SellExchangeDate string
	ProfitablePrice  decimal.Big
	Status           string
	Version          int32
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type OrderEvent struct {
	AggregateID uuid.UUID
	ParentID    uuid.UUID
	EventType   string
	Payload     []byte
	Version     int32
	CreatedAt   time.Time
}

type PickedStock struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	StockID   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type StakeConcentration struct {
	ID              uuid.UUID
	StockID         string
	ExchangeDate    string
	SumBuyShares    sql.NullInt64
	SumSellShares   sql.NullInt64
	AvgBuyPrice     decimal.Big
	AvgSellPrice    decimal.Big
	Concentration1  pgtype.Numeric
	Concentration5  pgtype.Numeric
	Concentration10 pgtype.Numeric
	Concentration20 pgtype.Numeric
	Concentration60 pgtype.Numeric
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       sql.NullTime
}

type Stock struct {
	ID        string
	Name      string
	Country   string
	Site      sql.NullString
	Category  sql.NullString
	Market    sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type ThreePrimary struct {
	ID                 uuid.UUID
	StockID            string
	ExchangeDate       string
	ForeignTradeShares sql.NullInt64
	TrustTradeShares   sql.NullInt64
	DealerTradeShares  sql.NullInt64
	HedgingTradeShares sql.NullInt64
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          sql.NullTime
}

type Transaction struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	OrderID      uuid.UUID
	OrderType    string
	CreditAmount decimal.Big
	DebitAmount  decimal.Big
	Status       string
	Version      int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TransactionEvent struct {
	AggregateID uuid.UUID
	ParentID    uuid.UUID
	EventType   string
	Payload     []byte
	Version     int32
	CreatedAt   time.Time
}

type User struct {
	ID               uuid.UUID
	FirstName        string
	LastName         string
	Email            string
	Phone            string
	Password         string
	SessionID        sql.NullString
	EmailConfirmedAt sql.NullTime
	PhoneConfirmedAt sql.NullTime
	CreatedAt        time.Time
	UpdatedAt        time.Time
	SessionExpiredAt sql.NullTime
	DeletedAt        sql.NullTime
}
