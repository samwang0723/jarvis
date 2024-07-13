package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	Time
	SessionExpiredAt *time.Time
	PhoneConfirmedAt *time.Time
	EmailConfirmedAt *time.Time
	FirstName        string
	LastName         string
	Email            string
	Phone            string
	Password         string
	SessionID        string
	ID
}

type UpdateSessionIDParams struct {
	SessionExpiredAt time.Time
	SessionID        string
	ID               uuid.UUID
}
