package entity

import "errors"

var (
	ErrProcessTimeout          = errors.New("process timeout")
	ErrStringIsNotHex          = errors.New("string is not a hex string")
	ErrTooMuchRequestToService = errors.New("too many requests to service")
	ErrInternalServer          = errors.New("internal server error")
)
