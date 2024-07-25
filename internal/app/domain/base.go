package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type ID struct {
	ID uuid.UUID
}

type Time struct {
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
