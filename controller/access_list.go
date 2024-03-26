package controller

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"isp-system-service/domain"
)

type AccessListService interface {
	GetById(ctx context.Context, appId int) ([]domain.MethodInfo, error)
	SetOne(ctx context.Context, request domain.AccessListSetOneRequest) (*domain.AccessListSetOneResponse, error)
	SetList(ctx context.Context, req domain.AccessListSetListRequest) ([]domain.MethodInfo, error)
}

type AccessList struct {
	service AccessListService
}

func NewAccessList(service AccessListService) AccessList {
	return AccessList{
		service: service,
	}
}

// GetById godoc
// @Tags accessList
// @Summary Получить список доступности методов для приложения
// @Description Возвращает список методов для приложения, для которых заданы настройки доступа
// @Accept  json
// @Produce  json
// @Param body body domain.Identity false "идентификатор приложения"
// @Success 200 {array} domain.MethodInfo "список доступности методов"
// @Failure 404 {object} domain.GrpcError
// @Failure 500 {object} domain.GrpcError
// @Router /access_list/get_by_id [POST]
func (c AccessList) GetById(ctx context.Context, req domain.Identity) ([]domain.MethodInfo, error) {
	result, err := c.service.GetById(ctx, req.Id)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, status.Errorf(codes.NotFound, "application not found")
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// SetOne godoc
// @Tags accessList
// @Summary Настроить доступность метода для приложения
// @Description Возвращает количество измененных строк
// @Accept  json
// @Produce  json
// @Param body body domain.AccessListSetOneRequest false "объект для настройки доступа"
// @Success 200 {object} domain.AccessListSetOneResponse "количество измененных строк"
// @Failure 404 {object} domain.GrpcError
// @Failure 500 {object} domain.GrpcError
// @Router /access_list/set_one [POST]
func (c AccessList) SetOne(ctx context.Context, req domain.AccessListSetOneRequest) (*domain.AccessListSetOneResponse, error) {
	result, err := c.service.SetOne(ctx, req)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, status.Errorf(codes.NotFound, "application not found")
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// SetList godoc
// @Tags accessList
// @Summary Настроить доступность списка методов для приложения
// @Description Возвращает список методов для приложения, для которых заданы настройки доступа
// @Accept  json
// @Produce  json
// @Param body body domain.AccessListSetListRequest false "объект настройки доступа"
// @Success 200 {array} domain.MethodInfo "список доступности методов"
// @Failure 404 {object} domain.GrpcError
// @Failure 500 {object} domain.GrpcError
// @Router /access_list/set_list [POST]
func (c AccessList) SetList(ctx context.Context, req domain.AccessListSetListRequest) ([]domain.MethodInfo, error) {
	result, err := c.service.SetList(ctx, req)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, status.Errorf(codes.NotFound, "application not found")
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}
