package service

import (
	"context"

	"github.com/pkg/errors"
	"isp-system-service/domain"
	"isp-system-service/entity"
)

type ServiceServiceRep interface {
	GetServiceById(ctx context.Context, id int) (*entity.Service, error)
	GetServiceByIdList(ctx context.Context, idList []int) ([]entity.Service, error)
	GetServiceByDomainId(ctx context.Context, domainIdList []int) ([]entity.Service, error)
	GetServiceByNameAndDomainId(ctx context.Context, name string, domainId int) (*entity.Service, error)
	CreateService(ctx context.Context, name string, desc string, domainId int) (*entity.Service, error)
	UpdateService(ctx context.Context, id int, name string, description string) (*entity.Service, error)
	DeleteService(ctx context.Context, idList []int) (int, error)
}

type ServiceDomainRep interface { // nolint:iface
	GetDomainById(ctx context.Context, id int) (*entity.Domain, error)
}

type Service struct {
	domainRep  ServiceDomainRep
	serviceRep ServiceServiceRep
}

func NewService(
	domainRep ServiceDomainRep,
	serviceRep ServiceServiceRep,
) Service {
	return Service{
		domainRep:  domainRep,
		serviceRep: serviceRep,
	}
}

func (s Service) GetById(ctx context.Context, id int) (*domain.Service, error) {
	serviceEntity, err := s.serviceRep.GetServiceById(ctx, id)
	if err != nil {
		return nil, errors.WithMessagef(err, "get service by id")
	}

	result := s.convertService(*serviceEntity)
	return &result, nil
}

func (s Service) GetByIdList(ctx context.Context, idList []int) ([]domain.Service, error) {
	serviceEntity, err := s.serviceRep.GetServiceByIdList(ctx, idList)
	if err != nil {
		return nil, errors.WithMessage(err, "get system by id list")
	}

	result := make([]domain.Service, len(serviceEntity))
	for i := range serviceEntity {
		result[i] = s.convertService(serviceEntity[i])
	}
	return result, nil
}

func (s Service) GetByDomainId(ctx context.Context, domainId int) ([]domain.Service, error) {
	serviceEntity, err := s.serviceRep.GetServiceByDomainId(ctx, []int{domainId})
	if err != nil {
		return nil, errors.WithMessage(err, "get system by id list")
	}

	result := make([]domain.Service, len(serviceEntity))
	for i := range serviceEntity {
		result[i] = s.convertService(serviceEntity[i])
	}
	return result, nil
}

func (s Service) CreateUpdate(ctx context.Context, req domain.ServiceCreateUpdateRequest) (*domain.Service, error) {
	req.DomainId = 1 // temporary use only 1 domain, soon domain entity will be removed

	existed, err := s.serviceRep.GetServiceByNameAndDomainId(ctx, req.Name, req.DomainId)
	switch {
	case errors.Is(err, domain.ErrServiceNotFound):
	case err != nil:
		return nil, errors.WithMessage(err, "get service by name and domain_id")
	}

	_, err = s.domainRep.GetDomainById(ctx, req.DomainId)
	if err != nil {
		return nil, errors.WithMessage(err, "get domain by id")
	}

	if req.Id == 0 {
		if existed != nil {
			return nil, domain.ErrDomainDuplicateName
		}

		serviceEntity, err := s.serviceRep.CreateService(ctx, req.Name, req.Description, req.DomainId)
		if err != nil {
			return nil, errors.WithMessage(err, "create service")
		}

		result := s.convertService(*serviceEntity)
		return &result, nil
	}

	if existed != nil && existed.Id != req.Id {
		return nil, domain.ErrDomainDuplicateName
	}

	_, err = s.serviceRep.GetServiceById(ctx, req.Id)
	if err != nil {
		return nil, errors.WithMessage(err, "get service by id")
	}

	serviceEntity, err := s.serviceRep.UpdateService(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		return nil, errors.WithMessage(err, "update service")
	}

	result := s.convertService(*serviceEntity)
	return &result, nil
}

func (s Service) Delete(ctx context.Context, idList []int) (int, error) {
	result, err := s.serviceRep.DeleteService(ctx, idList)
	if err != nil {
		return 0, errors.WithMessage(err, "delete service")
	}

	return result, nil
}

func (s Service) convertService(req entity.Service) domain.Service {
	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}
	result := domain.Service{
		Id:          req.Id,
		Name:        req.Name,
		Description: desc,
		DomainId:    req.DomainId,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}

	return result
}
