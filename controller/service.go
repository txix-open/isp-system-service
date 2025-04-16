package controller

import (
	"context"
	"fmt"

	"isp-system-service/domain"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/apierrors"
	"google.golang.org/grpc/codes"
)

type ServiceService interface {
	GetById(ctx context.Context, id int) (*domain.Service, error)
	GetByIdList(ctx context.Context, idList []int) ([]domain.Service, error)
	GetByDomainId(ctx context.Context, domainId int) ([]domain.Service, error)
	CreateUpdate(ctx context.Context, req domain.ServiceCreateUpdateRequest) (*domain.Service, error)
	Delete(ctx context.Context, idList []int) (int, error)
}

type Service struct {
	service ServiceService
}

func NewService(service ServiceService) Service {
	return Service{
		service: service,
	}
}

// GetById godoc
//
//	@Tags			service
//	@Summary		Получить сервис по идентификатору
//	@Description	Возвращает описание сервиса по его идентификатору
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.Identity	true	"Идентификатор сервиса"
//	@Success		200		{object}	domain.Service
//	@Failure		404		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/service/get_service_by_id [POST]
func (c Service) GetById(ctx context.Context, req domain.Identity) (*domain.Service, error) {
	result, err := c.service.GetById(ctx, req.Id)
	switch {
	case errors.Is(err, domain.ErrAppGroupNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeAppGroupNotFound,
			fmt.Sprintf("service with id %d not found", req.Id),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// Get godoc
//
//	@Tags			service
//	@Summary		Получить список сервисов
//	@Description	Возвращает список сервисов по их идентификаторам
//	@Accept			json
//	@Produce		json
//	@Param			body	body		[]integer	false	"Массив идентификаторов сервисов"
//	@Success		200		{array}		domain.Service
//	@Failure		500		{object}	apierrors.Error
//	@Router			/service/get_service [POST]
func (c Service) Get(ctx context.Context, req []int) ([]domain.Service, error) {
	return c.service.GetByIdList(ctx, req)
}

// GetByDomainId godoc
//
//	@Tags			service
//	@Summary		Получить список сервисов по идентификатору домена
//	@Description	Возвращает список сервисов по идентификатору домена
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.Identity	true	"Идентификатор домена"
//	@Success		200		{array}		domain.Service
//	@Failure		500		{object}	apierrors.Error
//	@Router			/service/get_services_by_domain_id [POST]
func (c Service) GetByDomainId(ctx context.Context, req domain.Identity) ([]domain.Service, error) {
	return c.service.GetByDomainId(ctx, req.Id)
}

// CreateUpdate godoc
//
//	@Tags			service
//	@Summary		Создать/обновить сервис
//	@Description	Если сервис с такими идентификатором существует, то обновляет данные, если нет, то добавляет данные в базу
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.ServiceCreateUpdateRequest	true	"Объект сервиса"
//	@Success		200		{object}	domain.Service
//	@Failure		400		{object}	apierrors.Error
//	@Failure		404		{object}	apierrors.Error
//	@Failure		409		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/service/create_update_service [POST]
func (c Service) CreateUpdate(ctx context.Context, req domain.ServiceCreateUpdateRequest) (*domain.Service, error) {
	result, err := c.service.CreateUpdate(ctx, req)
	switch {
	case errors.Is(err, domain.ErrDomainNotFound):
		return nil, apierrors.NewBusinessError(
			domain.ErrCodeDomainNotFound,
			fmt.Sprintf("domain with id %d not found", req.DomainId),
			err,
		)
	case errors.Is(err, domain.ErrAppGroupNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeAppGroupNotFound,
			fmt.Sprintf("service with id %d not found", req.Id),
			err,
		)
	case errors.Is(err, domain.ErrAppGroupDuplicateName):
		return nil, apierrors.New(
			codes.AlreadyExists,
			domain.ErrCodeAppGroupNotFound,
			fmt.Sprintf("service with name %s already exists", req.Name),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// Delete godoc
//
//	@Tags			service
//	@Summary		Удалить сервисы
//	@Description	Удаляет сервисов по списку их идентификаторов, возвращает количество удаленных сервисов
//	@Accept			json
//	@Produce		json
//	@Param			body	body		[]integer	true	"Массив идентификаторов сервисов"
//	@Success		200		{object}	domain.DeleteResponse
//	@Failure		400		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/service/delete_service [POST]
func (c Service) Delete(ctx context.Context, req []int) (*domain.DeleteResponse, error) {
	if len(req) == 0 {
		return nil, apierrors.NewBusinessError(domain.ErrCodeInvalidRequest,
			"At least one id are required", errors.New("invalid id count"))
	}

	result, err := c.service.Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	return &domain.DeleteResponse{
		Deleted: result,
	}, nil
}
