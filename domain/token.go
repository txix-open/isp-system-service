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
	AppId  int `valid:"required~Required"`
	Tokens []string
}

type TokenCreateRequest struct {
	AppId        int `valid:"required~Required"`
	ExpireTimeMs int `valid:"required~Required"`
}
