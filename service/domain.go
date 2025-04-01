package service

import (
	"context"

	"github.com/pkg/errors"
	"isp-system-service/domain"
	"isp-system-service/entity"
)

type Domain struct {
	repo DomainRepo
}

func NewDomain(repo DomainRepo) Domain {
	return Domain{
		repo: repo,
	}
}

func (s Domain) GetById(ctx context.Context, id int) (*domain.Domain, error) {
	domainEntity, err := s.repo.GetDomainById(ctx, id)
	if err != nil {
		return nil, errors.WithMessage(err, "get domain by id")
	}

	result := s.convertDomain(*domainEntity)
	return &result, nil
}

func (s Domain) GetBySystemId(ctx context.Context, systemId int) ([]domain.Domain, error) {
	domainEntity, err := s.repo.GetDomainBySystemId(ctx, systemId)
	if err != nil {
		return nil, errors.WithMessage(err, "get domain by system_id")
	}

	result := make([]domain.Domain, len(domainEntity))
	for i := range domainEntity {
		result[i] = s.convertDomain(domainEntity[i])
	}
	return result, nil
}

func (s Domain) CreateUpdate(ctx context.Context, req domain.DomainCreateUpdateRequest, systemId int) (*domain.Domain, error) {
	existed, err := s.repo.GetDomainByNameAndSystemId(ctx, req.Name, systemId)
	switch {
	case errors.Is(err, domain.ErrDomainNotFound):
	case err != nil:
		return nil, errors.WithMessage(err, "get domain by name and system_id")
	}

	if req.Id == 0 {
		if existed != nil {
			return nil, domain.ErrDomainDuplicateName
		}

		domainEntity, err := s.repo.CreateDomain(ctx, req.Name, req.Description, systemId)
		if err != nil {
			return nil, errors.WithMessage(err, "create domain")
		}

		result := s.convertDomain(*domainEntity)
		return &result, nil
	}

	if existed != nil && existed.Id != req.Id {
		return nil, domain.ErrDomainDuplicateName
	}

	_, err = s.repo.GetDomainById(ctx, req.Id)
	if err != nil {
		return nil, errors.WithMessage(err, "get domain by id")
	}

	domainEntity, err := s.repo.UpdateDomain(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		return nil, errors.WithMessage(err, "update domain")
	}

	result := s.convertDomain(*domainEntity)
	return &result, nil
}

func (s Domain) Delete(ctx context.Context, idList []int) (int, error) {
	result, err := s.repo.DeleteDomain(ctx, idList)
	if err != nil {
		return 0, errors.WithMessage(err, "delete domain")
	}

	return result, nil
}

func (s Domain) convertDomain(req entity.Domain) domain.Domain {
	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}
	result := domain.Domain{
		Id:          req.Id,
		Name:        req.Name,
		Description: desc,
		SystemId:    req.SystemId,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}

	return result
}
