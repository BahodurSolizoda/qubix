package model

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrAttemptCompleted   = errors.New("attempt is already completed")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrDuplicate          = errors.New("login or email already exists")
)
