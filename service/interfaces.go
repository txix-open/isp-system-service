package service

import (
	"context"
	"isp-system-service/entity"
)

type TokenRepo interface {
	GetTokenByAppIdList(ctx context.Context, appIdList []int) ([]entity.Token, error)
	GetTokenById(ctx context.Context, token string) (*entity.Token, error)
}

type ApplicationRepo interface {
	GetApplicationById(ctx context.Context, id int) (*entity.Application, error)
	GetApplicationByIdList(ctx context.Context, idList []int) ([]entity.Application, error)
	GetApplicationByAppGroupIdList(ctx context.Context, appGroupIdList []int) ([]entity.Application, error)
	CreateApplication(ctx context.Context, id int, name string, desc string, appGroupId int, appType string) (*entity.Application, error)
	UpdateApplication(ctx context.Context, id int, name string, description string) (*entity.Application, error)
	UpdateApplicationWithNewId(ctx context.Context, oldId int, newId int, name string, description string) (*entity.Application, error)
	NextApplicationId(ctx context.Context) (int, error)
	GetAllApplications(ctx context.Context) ([]entity.Application, error)
}

type AppGroupRepo interface {
	GetAppGroupById(ctx context.Context, id int) (*entity.AppGroup, error)
	GetAppGroupByIdList(ctx context.Context, idList []int) ([]entity.AppGroup, error)
	GetAllAppGroups(ctx context.Context) ([]entity.AppGroup, error)
	GetAppGroupByName(ctx context.Context, name string) (*entity.AppGroup, error)
	CreateAppGroup(ctx context.Context, name string, desc string) (*entity.AppGroup, error)
	UpdateAppGroup(ctx context.Context, id int, name string, description string) (*entity.AppGroup, error)
	DeleteAppGroup(ctx context.Context, idList []int) (int, error)
}
