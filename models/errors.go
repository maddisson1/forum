package models

import "errors"

type CustomError struct {
	ErrorCode int
	ErrorMsg  string
}

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
