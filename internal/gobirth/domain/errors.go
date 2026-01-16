package domain

import "errors"

var (
	ErrMissingPhone = errors.New("missing phone number")
	ErrInvalidPhone = errors.New("invalid phone number")
	ErrMissingName  = errors.New("missing contact name")
)
