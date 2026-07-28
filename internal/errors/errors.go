package errors

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrNotFound       = errors.New("resource not found")
	ErrInternal       = errors.New("internal server error")
)