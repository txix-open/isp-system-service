package controller

import (
	"isp-system-service/entity"
)

type DeleteResponse struct {
	Deleted int
}

type CountResponse struct {
	Count int
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

type ModuleMethods map[string][]MethodInfo

type MethodInfo struct {
	Method string
	Value  bool
}

type SetListRequest struct {
	AppId   int32
	Methods []MethodInfo
}
