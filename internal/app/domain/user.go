package domain

import (
	"time"
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
	Time
}
