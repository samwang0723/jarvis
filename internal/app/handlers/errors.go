package handlers

import "errors"

var (
	errOrderTypeNotAllowed = errors.New("order type not allowed")
	errInvalidCaptcha      = errors.New("invalid captcha")
)
