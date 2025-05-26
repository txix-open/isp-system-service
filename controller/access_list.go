package controller

import (
	"context"
	"fmt"

	"isp-system-service/domain"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/apierrors"
	"google.golang.org/grpc/codes"
)

type AccessListService interface {
	GetById(ctx context.Context, appId int) ([]domain.MethodInfo, error)
	SetOne(ctx context.Context, request domain.AccessListSetOneRequest) (*domain.AccessListSetOneResponse, error)
	SetList(ctx context.Context, req domain.AccessListSetListRequest) ([]domain.MethodInfo, error)
	DeleteList(ctx context.Context, req domain.AccessListDeleteListRequest) error
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
//
//	@Tags			accessList
//	@Summary		Получить список доступности методов для приложения
//	@Description	Возвращает список методов для приложения, для которых заданы настройки доступа
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.Identity		false	"идентификатор приложения"
//	@Success		200		{array}		domain.MethodInfo	"список доступности методов"
//	@Failure		404		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/access_list/get_by_id [POST]
func (c AccessList) GetById(ctx context.Context, req domain.Identity) ([]domain.MethodInfo, error) {
	result, err := c.service.GetById(ctx, req.Id)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeApplicationNotFound,
			fmt.Sprintf("application with id %d not found", req.Id),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// SetOne godoc
//
//	@Tags			accessList
//	@Summary		Настроить доступность метода для приложения
//	@Description	Возвращает количество измененных строк
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.AccessListSetOneRequest	false	"объект для настройки доступа"
//	@Success		200		{object}	domain.AccessListSetOneResponse	"количество измененных строк"
//	@Failure		404		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/access_list/set_one [POST]
func (c AccessList) SetOne(ctx context.Context, req domain.AccessListSetOneRequest) (*domain.AccessListSetOneResponse, error) {
	result, err := c.service.SetOne(ctx, req)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeApplicationNotFound,
			fmt.Sprintf("application with id %d not found", req.AppId),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// SetList godoc
//
//	@Tags			accessList
//	@Summary		Настроить доступность списка методов для приложения
//	@Description	Возвращает список методов для приложения, для которых заданы настройки доступа
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.AccessListSetListRequest	false	"объект настройки доступа"
//	@Success		200		{array}		domain.MethodInfo				"список доступности методов"
//	@Failure		404		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/access_list/set_list [POST]
func (c AccessList) SetList(ctx context.Context, req domain.AccessListSetListRequest) ([]domain.MethodInfo, error) {
	result, err := c.service.SetList(ctx, req)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeApplicationNotFound,
			fmt.Sprintf("application with id %d not found", req.AppId),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// DeleteList godoc
//
//	@Tags			accessList
//	@Summary		Удалить список доступных методов для приложения
//	@Description	Удаляет заданный список методов для приложения
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.AccessListDeleteListRequest	false	"тело запроса"
//	@Success		200		{object}	any
//	@Failure		404		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/access_list/delete_list [POST]
func (c AccessList) DeleteList(ctx context.Context, req domain.AccessListDeleteListRequest) error {
	err := c.service.DeleteList(ctx, req)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return apierrors.New(
			codes.NotFound,
			domain.ErrCodeApplicationNotFound,
			fmt.Sprintf("application with id %d not found", req.AppId),
			err,
		)
	default:
		return err
	}
}
