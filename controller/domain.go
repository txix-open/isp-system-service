package controller

import (
	"context"
	"fmt"

	"isp-system-service/domain"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/apierrors"
	"google.golang.org/grpc/codes"
)

type DomainService interface {
	GetById(ctx context.Context, id int) (*domain.Domain, error)
	GetBySystemId(ctx context.Context, systemId int) ([]domain.Domain, error)
	CreateUpdate(ctx context.Context, req domain.DomainCreateUpdateRequest, systemId int) (*domain.Domain, error)
	Delete(ctx context.Context, idList []int) (int, error)
}

type Domain struct {
	service DomainService
}

func NewDomain(service DomainService) Domain {
	return Domain{
		service: service,
	}
}

// GetById godoc
//
//	@Tags			domain
//	@Summary		Получить домен по идентификатору
//	@Description	Возвращает описание домена по его идентификатору
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.Identity	true	"Идентификатор домена"
//	@Success		200		{object}	domain.Domain
//	@Failure		404		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/domain/get_domain_by_id [POST]
func (c Domain) GetById(ctx context.Context, req domain.Identity) (*domain.Domain, error) {
	result, err := c.service.GetById(ctx, req.Id)
	switch {
	case errors.Is(err, domain.ErrDomainNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeDomainNotFound,
			fmt.Sprintf("domain with id %d not found", req.Id),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// GetBySystemId godoc
//
//	@Tags			domain
//	@Summary		Получить домены по идентификатору системы
//	@Description	Возвращает список доменов по системному идентификатору
//	@Accept			json
//	@Produce		json
//	@Param			body	body		integer	false	"Идентификатор системы"
//	@Success		200		{array}		domain.Domain
//	@Failure		500		{object}	apierrors.Error
//	@Router			/domain/get_domains_by_system_id [POST]
func (c Domain) GetBySystemId(ctx context.Context) ([]domain.Domain, error) {
	return c.service.GetBySystemId(ctx, domain.DefaultSystemId)
}

// CreateUpdate godoc
//
//	@Tags			domain
//	@Summary		Создать/обновить домен
//	@Description	Если домен с такими идентификатором существует, то обновляет данные, если нет, то добавляет данные в базу
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.DomainCreateUpdateRequest	true	"Объект домена"
//	@Success		200		{object}	domain.Domain
//	@Failure		500		{object}	apierrors.Error
//	@Router			/domain/create_update_domain [POST]
func (c Domain) CreateUpdate(ctx context.Context, req domain.DomainCreateUpdateRequest) (*domain.Domain, error) {
	result, err := c.service.CreateUpdate(ctx, req, domain.DefaultSystemId)
	switch {
	case errors.Is(err, domain.ErrSystemNotFound):
		return nil, apierrors.NewBusinessError(
			domain.ErrCodeSystemNotFound,
			fmt.Sprintf("system with id %d not found", domain.DefaultSystemId),
			err,
		)
	case errors.Is(err, domain.ErrDomainNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeDomainNotFound,
			fmt.Sprintf("domain with id %d not found", req.Id),
			err,
		)
	case errors.Is(err, domain.ErrDomainDuplicateName):
		return nil, apierrors.New(
			codes.AlreadyExists,
			domain.ErrCodeDomainDuplicateName,
			fmt.Sprintf("domain with name %s already exists", req.Name),
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
//	@Tags			domain
//	@Summary		Удаление доменов
//	@Description	Удаляет домены по списку их идентификаторов, возвращает количество удаленных доменов
//	@Accept			json
//	@Produce		json
//	@Param			body	body		[]integer	false	"Массив идентификаторов доменов"
//	@Success		200		{object}	domain.DeleteResponse
//	@Failure		400		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/domain/delete_domains [POST]
func (c Domain) Delete(ctx context.Context, req []int) (*domain.DeleteResponse, error) {
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
