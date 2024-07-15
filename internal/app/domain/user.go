package domain

import (
	"regexp"
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

func (u *User) Validate() error {
	if u.Email == "" {
		return &DataMissingError{dataType: "email"}
	}
	if !isValidEmail(u.Email) {
		return &DataValidationError{dataType: "email"}
	}

	if u.Phone == "" {
		return &DataMissingError{dataType: "phone"}
	}
	if !isValidPhone(u.Phone) {
		return &DataValidationError{dataType: "phone"}
	}

	return nil
}

func isValidEmail(email string) bool {
	// Simple regex for email validation
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func isValidPhone(phone string) bool {
	// Simple regex for phone validation
	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return re.MatchString(phone)
}
