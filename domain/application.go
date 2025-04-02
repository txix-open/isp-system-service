package domain

import (
	"time"
)

const (
	ApplicationSystemType = "SYSTEM"
	ApplicationMobileType = "MOBILE"
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

type CreateApplicationRequest struct {
	Id                 int
	Name               string `validate:"required"`
	Description        string
	ApplicationGroupId int    `validate:"required"`
	Type               string `validate:"required,oneof=SYSTEM MOBILE"`
}
