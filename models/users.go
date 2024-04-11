package models

import (
	"forum/pkg/validator"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	hashedPassword string
	Created        time.Time
}

type UserLoginForm struct {
	Email               string
	Password            string
	validator.Validator `form:"-"`
}

type UserSignupForm struct {
	Name                string
	Email               string
	Password            string
	validator.Validator `form:"-"`
}
