package entity

import (
	"time"
)

type Token struct {
	Token      string
	AppId      int
	ExpireTime int
	CreatedAt  time.Time
}
