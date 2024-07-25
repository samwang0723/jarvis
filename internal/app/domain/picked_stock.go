package domain

import "github.com/gofrs/uuid/v5"

type PickedStock struct {
	Time
	StockID string
	ID
	UserID uuid.UUID
}
