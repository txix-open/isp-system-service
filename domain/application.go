package domain

import (
	"time"
)

type Application struct {
	Id          int
	Name        string
	Description string
	ServiceId   int
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ApplicationCreateUpdateRequest struct {
	Id          int
	Name        string `validate:"required"`
	Description string
	ServiceId   int    `validate:"required"`
	Type        string `validate:"required,oneof=SYSTEM MOBILE"`
}

type ApplicationWithTokens struct {
	App    Application
	Tokens []Token
}

type ApplicationSimple struct {
	Id          int
	Name        string
	Description string
	Type        string
	Tokens      []Token
}
