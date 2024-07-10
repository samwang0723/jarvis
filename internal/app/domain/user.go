package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID
	FirstName        string
	LastName         string
	Email            string
	Phone            string
	Password         string
	SessionID        string
	SessionExpiredAt *time.Time
	PhoneConfirmedAt *time.Time
	EmailConfirmedAt *time.Time
	Time
}

type UpdateSessionIDParams struct {
	SessionID        string
	SessionExpiredAt time.Time
	ID               uuid.UUID
}
