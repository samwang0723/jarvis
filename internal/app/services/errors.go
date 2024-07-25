package services

import "errors"

var (
	errCannotCastDailyClose      = errors.New("cannot cast interface to *dto.DailyClose")
	errUnableToChainTransactions = errors.New("unable to create chain transactions")
	errUserNotFound              = errors.New("user not found")
	errUserPasswordNotMatch      = errors.New("user login credential not match")
)
