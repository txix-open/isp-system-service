package controller

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"isp-system-service/domain"
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
// @Tags domain
// @Summary Получить домен по идентификатору
// @Description Возвращает описание домена по его идентификатору
// @Accept  json
// @Produce  json
// @Param body body domain.Identity true "Идентификатор домена"
// @Success 200 {object} domain.Domain
// @Failure 404 {object} domain.GrpcError
// @Failure 500 {object} domain.GrpcError
// @Router /domain/get_domain_by_id [POST]
func (c Domain) GetById(ctx context.Context, req domain.Identity) (*domain.Domain, error) {
	result, err := c.service.GetById(ctx, req.Id)
	switch {
	case errors.Is(err, domain.ErrDomainNotFound):
		return nil, status.Errorf(codes.NotFound, "domain with id %d not found", req.Id)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// GetBySystemId godoc
// @Tags domain
// @Summary Получить домены по идентификатору системы
// @Description Возвращает список доменов по системному идентификатору
// @Accept  json
// @Produce  json
// @Param body body integer false "Идентификатор системы"
// @Success 200 {array} domain.Domain
// @Failure 500 {object} domain.GrpcError
// @Router /domain/get_domains_by_system_id [POST]
func (c Domain) GetBySystemId(ctx context.Context) ([]domain.Domain, error) {
	return c.service.GetBySystemId(ctx, domain.DefaultSystemId)
}

// CreateUpdate godoc
// @Tags domain
// @Summary Создать/обновить домен
// @Description Если домен с такими идентификатором существует, то обновляет данные, если нет, то добавляет данные в базу
// @Accept  json
// @Produce  json
// @Param body body domain.DomainCreateUpdateRequest true "Объект домена"
// @Success 200 {object} domain.Domain
// @Failure 500 {object} domain.GrpcError
// @Router /domain/create_update_domain [POST]
func (c Domain) CreateUpdate(ctx context.Context, req domain.DomainCreateUpdateRequest) (*domain.Domain, error) {
	result, err := c.service.CreateUpdate(ctx, req, domain.DefaultSystemId)
	switch {
	case errors.Is(err, domain.ErrSystemNotFound):
		return nil, status.Errorf(codes.InvalidArgument, "system with id %d not found", domain.DefaultSystemId)
	case errors.Is(err, domain.ErrDomainNotFound):
		return nil, status.Errorf(codes.NotFound, "domain with id %d not found", req.Id)
	case errors.Is(err, domain.ErrDomainDuplicateName):
		return nil, status.Errorf(codes.AlreadyExists, "domain with name %s already exists", req.Name)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// Delete godoc
// @Tags domain
// @Summary Удаление доменов
// @Description Удаляет домены по списку их идентификаторов, возвращает количество удаленных доменов
// @Accept  json
// @Produce  json
// @Param body body []integer false "Массив идентификаторов доменов"
// @Success 200 {object} domain.DeleteResponse
// @Failure 400 {object} domain.GrpcError
// @Failure 500 {object} domain.GrpcError
// @Router /domain/delete_domains [POST]
func (c Domain) Delete(ctx context.Context, req []int) (*domain.DeleteResponse, error) {
	if len(req) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "at least one id are required")
	}

	result, err := c.service.Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	return &domain.DeleteResponse{
		Deleted: result,
	}, nil
}
