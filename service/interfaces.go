package service

import (
	"context"
	"isp-system-service/entity"
)

type TokenRepo interface {
	GetTokenByAppIdList(ctx context.Context, appIdList []int) ([]entity.Token, error)
}

type DomainRepo interface {
	GetDomainById(ctx context.Context, id int) (*entity.Domain, error)
	GetDomainBySystemId(ctx context.Context, systemId int) ([]entity.Domain, error)
	GetDomainByNameAndSystemId(ctx context.Context, name string, systemId int) (*entity.Domain, error)
	CreateDomain(ctx context.Context, name string, desc string, systemId int) (*entity.Domain, error)
	UpdateDomain(ctx context.Context, id int, name string, description string) (*entity.Domain, error)
	DeleteDomain(ctx context.Context, idList []int) (int, error)
}

type ApplicationRepo interface {
	GetApplicationById(ctx context.Context, id int) (*entity.Application, error)
	GetApplicationByIdList(ctx context.Context, idList []int) ([]entity.Application, error)
	GetApplicationByServiceIdList(ctx context.Context, serviceIdList []int) ([]entity.Application, error)
	GetApplicationByNameAndServiceId(ctx context.Context, name string, serviceId int) (*entity.Application, error)
	CreateApplication(ctx context.Context, name string, desc string, serviceId int, appType string) (*entity.Application, error)
	UpdateApplication(ctx context.Context, id int, name string, description string) (*entity.Application, error)
}

type ServiceRepo interface {
	GetServiceById(ctx context.Context, id int) (*entity.Service, error)
	GetServiceByIdList(ctx context.Context, idList []int) ([]entity.Service, error)
	GetServiceByDomainId(ctx context.Context, domainIdList []int) ([]entity.Service, error)
	GetServiceByNameAndDomainId(ctx context.Context, name string, domainId int) (*entity.Service, error)
	CreateService(ctx context.Context, name string, desc string, domainId int) (*entity.Service, error)
	UpdateService(ctx context.Context, id int, name string, description string) (*entity.Service, error)
	DeleteService(ctx context.Context, idList []int) (int, error)
}
