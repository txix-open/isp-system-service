package domain

import (
	"time"
)

type Application struct {
	Id                 int
	Name               string
	Description        string
	ApplicationGroupId int
	Type               string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ApplicationCreateUpdateRequest struct {
	Id                 int
	Name               string `valid:"required~Required"`
	Description        string
	ApplicationGroupId int    `valid:"required~Required"`
	Type               string `valid:"required~Required,in(SYSTEM|MOBILE)"`
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
