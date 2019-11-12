package domain

import "isp-system-service/entity"

type (
	DeleteResponse struct {
		Deleted int
	}

	CountResponse struct {
		Count int
	}

	AppWithToken struct {
		App    entity.Application
		Tokens []entity.Token
	}

	DomainWithServices struct {
		Id          int32
		Name        string
		Description string
		Services    []*ServiceWithApps
	}

	ServiceWithApps struct {
		Id          int32
		Name        string
		Description string
		Apps        []*SimpleApp
	}

	SimpleApp struct {
		Id          int32
		Name        string
		Description string
		Type        string
		Tokens      []entity.Token
	}

	MethodInfo struct {
		Method string
		Value  bool
	}

	SetListRequest struct {
		AppId   int32
		Methods []MethodInfo
	}
)
