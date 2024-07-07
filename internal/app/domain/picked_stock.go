package domain

import "github.com/gofrs/uuid/v5"

type PickedStock struct {
	ID
	UserID  uuid.UUID
	StockID string
	Time
}
