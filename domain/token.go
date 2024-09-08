package domain

import (
	"time"
)

type Token struct {
	Token      string
	AppId      int
	ExpireTime int
	CreatedAt  time.Time
}

type TokenRevokeRequest struct {
	AppId  int `validate:"required"`
	Tokens []string
}

type TokenCreateRequest struct {
	AppId        int `validate:"required"`
	ExpireTimeMs int `validate:"required"`
}
