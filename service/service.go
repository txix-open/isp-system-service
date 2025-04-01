package service

import (
	"context"

	"isp-system-service/domain"
	"isp-system-service/entity"

	"github.com/pkg/errors"
)

type Service struct {
	domainRepo  DomainRepo
	serviceRepo AppGroupRepo
}

func NewService(
	domainRepo DomainRepo,
	serviceRepo AppGroupRepo,
) Service {
	return Service{
		domainRepo:  domainRepo,
		serviceRepo: serviceRepo,
	}
}

func (s Service) GetById(ctx context.Context, id int) (*domain.Service, error) {
	serviceEntity, err := s.serviceRepo.GetAppGroupById(ctx, id)
	if err != nil {
		return nil, errors.WithMessage(err, "get service by id")
	}

	result := s.convertService(*serviceEntity)
	return &result, nil
}

func (s Service) GetByIdList(ctx context.Context, idList []int) ([]domain.Service, error) {
	serviceEntity, err := s.serviceRepo.GetAppGroupByIdList(ctx, idList)
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
	serviceEntity, err := s.serviceRepo.GetAppGroupByDomainId(ctx, []int{domainId})
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

	existed, err := s.serviceRepo.GetAppGroupByNameAndDomainId(ctx, req.Name, req.DomainId)
	switch {
	case errors.Is(err, domain.ErrAppGroupNotFound):
	case err != nil:
		return nil, errors.WithMessage(err, "get service by name and domain_id")
	}

	_, err = s.domainRepo.GetDomainById(ctx, req.DomainId)
	if err != nil {
		return nil, errors.WithMessage(err, "get domain by id")
	}

	if req.Id == 0 {
		if existed != nil {
			return nil, domain.ErrDomainDuplicateName
		}

		serviceEntity, err := s.serviceRepo.CreateAppGroup(ctx, req.Name, req.Description, req.DomainId)
		if err != nil {
			return nil, errors.WithMessage(err, "create service")
		}

		result := s.convertService(*serviceEntity)
		return &result, nil
	}

	if existed != nil && existed.Id != req.Id {
		return nil, domain.ErrDomainDuplicateName
	}

	_, err = s.serviceRepo.GetAppGroupById(ctx, req.Id)
	if err != nil {
		return nil, errors.WithMessage(err, "get service by id")
	}

	serviceEntity, err := s.serviceRepo.UpdateAppGroup(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		return nil, errors.WithMessage(err, "update service")
	}

	result := s.convertService(*serviceEntity)
	return &result, nil
}

func (s Service) Delete(ctx context.Context, idList []int) (int, error) {
	result, err := s.serviceRepo.DeleteAppGroup(ctx, idList)
	if err != nil {
		return 0, errors.WithMessage(err, "delete service")
	}

	return result, nil
}

func (s Service) convertService(req entity.AppGroup) domain.Service {
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
