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
	GetApplicationByAppGroupIdList(ctx context.Context, appGroupIdList []int) ([]entity.Application, error)
	GetApplicationByNameAndAppGroupId(ctx context.Context, name string, appGroupId int) (*entity.Application, error)
	CreateApplication(ctx context.Context, id int, name string, desc string, appGroupId int, appType string) (*entity.Application, error)
	UpdateApplication(ctx context.Context, id int, name string, description string) (*entity.Application, error)
	UpdateApplicationWithNewId(ctx context.Context, oldId int, newId int, name string, description string) (*entity.Application, error)
	NextApplicationId(ctx context.Context) (int, error)
	GetAllApplications(ctx context.Context) ([]entity.Application, error)
}

type AppGroupRepo interface {
	GetAppGroupById(ctx context.Context, id int) (*entity.AppGroup, error)
	GetAppGroupByIdList(ctx context.Context, idList []int) ([]entity.AppGroup, error)
	GetAppGroupByDomainId(ctx context.Context, domainIdList []int) ([]entity.AppGroup, error)
	GetAllAppGroups(ctx context.Context) ([]entity.AppGroup, error)
	GetAppGroupByNameAndDomainId(ctx context.Context, name string, domainId int) (*entity.AppGroup, error)
	CreateAppGroup(ctx context.Context, name string, desc string, domainId int) (*entity.AppGroup, error)
	UpdateAppGroup(ctx context.Context, id int, name string, description string) (*entity.AppGroup, error)
	DeleteAppGroup(ctx context.Context, idList []int) (int, error)
}
