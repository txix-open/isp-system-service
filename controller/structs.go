package controller

import (
	"isp-system-service/entity"
)

type DeleteResponse struct {
	Deleted int
}

type Identity struct {
	Id int32 `json:"id" valid:"required~Required"`
}

type AppWithToken struct {
	App    entity.Application
	Tokens []entity.Token
}

type SimpleApp struct {
	Id          int32
	Name        string
	Description string
	Type        string
	Tokens      []entity.Token
}

type ServiceWithApps struct {
	Id          int32
	Name        string
	Description string
	Apps        []*SimpleApp
}

type DomainWithServices struct {
	Id          int32
	Name        string
	Description string
	Services    []*ServiceWithApps
}

type RevokeTokensRequest struct {
	AppId  int32 `valid:"required~Required"`
	Tokens []string
}
type CreateTokenRequest struct {
	AppId        int32 `valid:"required~Required"`
	ExpireTimeMs int64
}
