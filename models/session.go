package models

import "time"

type Session struct {
	UserID       int
	SessionToken string
	ExpiryAt     time.Time
}
