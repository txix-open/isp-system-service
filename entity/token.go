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

type AuthData struct {
	SystemId           int
	DomainId           int
	ApplicationGroupId int
	AppId              int
	ExpireTime         int
	CreatedAt          time.Time
}
