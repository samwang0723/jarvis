package services

import "errors"

var (
	errCannotCastDailyClose      = errors.New("cannot cast interface to *dto.DailyClose")
	errUnableToChainTransactions = errors.New("unable to create chain transactions")
	errInvalidJWTToken           = errors.New("invalid jwt token")
)
