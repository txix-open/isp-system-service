package service

import (
	"context"

	"github.com/pkg/errors"
	"isp-system-service/domain"
	"isp-system-service/entity"
)

type IApplicationGroupRep interface {
	GetApplicationGroupById(ctx context.Context, id int) (*entity.ApplicationGroup, error)
	GetApplicationGroupByIdList(ctx context.Context, idList []int) ([]entity.ApplicationGroup, error)
	GetApplicationGroupByDomainId(ctx context.Context, domainIdList []int) ([]entity.ApplicationGroup, error)
	GetSApplicationGroupByNameAndDomainId(ctx context.Context, name string, domainId int) (*entity.ApplicationGroup, error)
	CreateApplicationGroup(ctx context.Context, name string, desc string, domainId int) (*entity.ApplicationGroup, error)
	UpdateApplicationGroup(ctx context.Context, id int, name string, description string) (*entity.ApplicationGroup, error)
	DeleteApplicationGroup(ctx context.Context, idList []int) (int, error)
}

type IApplicationGroupDomainRep interface {
	GetDomainById(ctx context.Context, id int) (*entity.Domain, error)
}

type ApplicationGroup struct {
	domainRep IApplicationGroupDomainRep
	groupRep  IApplicationGroupRep
}

func NewApplicationGroup(
	domainRep IApplicationGroupDomainRep,
	groupRep IApplicationGroupRep,
) ApplicationGroup {
	return ApplicationGroup{
		domainRep: domainRep,
		groupRep:  groupRep,
	}
}

func (s ApplicationGroup) GetById(ctx context.Context, id int) (*domain.ApplicationGroup, error) {
	applicationGroupEntity, err := s.groupRep.GetApplicationGroupById(ctx, id)
	if err != nil {
		return nil, errors.WithMessagef(err, "get application group by id")
	}

	result := s.convertApplicationGroup(*applicationGroupEntity)
	return &result, nil
}

func (s ApplicationGroup) GetByIdList(ctx context.Context, idList []int) ([]domain.ApplicationGroup, error) {
	applicationGroupEntity, err := s.groupRep.GetApplicationGroupByIdList(ctx, idList)
	if err != nil {
		return nil, errors.WithMessage(err, "get system by id list")
	}

	result := make([]domain.ApplicationGroup, len(applicationGroupEntity))
	for i := range applicationGroupEntity {
		result[i] = s.convertApplicationGroup(applicationGroupEntity[i])
	}
	return result, nil
}

func (s ApplicationGroup) GetByDomainId(ctx context.Context, domainId int) ([]domain.ApplicationGroup, error) {
	applicationGroupEntity, err := s.groupRep.GetApplicationGroupByDomainId(ctx, []int{domainId})
	if err != nil {
		return nil, errors.WithMessage(err, "get system by id list")
	}

	result := make([]domain.ApplicationGroup, len(applicationGroupEntity))
	for i := range applicationGroupEntity {
		result[i] = s.convertApplicationGroup(applicationGroupEntity[i])
	}
	return result, nil
}

func (s ApplicationGroup) CreateUpdate(ctx context.Context, req domain.ApplicationGroupCreateUpdateRequest) (*domain.ApplicationGroup, error) {
	req.DomainId = 1 // Хардкодом проставляем domainId = 1

	existed, err := s.groupRep.GetSApplicationGroupByNameAndDomainId(ctx, req.Name, req.DomainId)
	switch {
	case errors.Is(err, domain.ErrApplicationGroupNotFound):
	case err != nil:
		return nil, errors.WithMessage(err, "get application group by name and domain_id")
	}

	_, err = s.domainRep.GetDomainById(ctx, req.DomainId)
	if err != nil {
		return nil, errors.WithMessage(err, "get domain by id")
	}

	if req.Id == 0 {
		if existed != nil {
			return nil, domain.ErrDomainDuplicateName
		}

		applicationGroupEntity, err := s.groupRep.CreateApplicationGroup(ctx, req.Name, req.Description, req.DomainId)
		if err != nil {
			return nil, errors.WithMessage(err, "create application group")
		}

		result := s.convertApplicationGroup(*applicationGroupEntity)
		return &result, nil
	}

	if existed != nil && existed.Id != req.Id {
		return nil, domain.ErrDomainDuplicateName
	}

	_, err = s.groupRep.GetApplicationGroupById(ctx, req.Id)
	if err != nil {
		return nil, errors.WithMessage(err, "get application group by id")
	}

	applicationGroupEntity, err := s.groupRep.UpdateApplicationGroup(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		return nil, errors.WithMessage(err, "update application group")
	}

	result := s.convertApplicationGroup(*applicationGroupEntity)
	return &result, nil
}

func (s ApplicationGroup) Delete(ctx context.Context, idList []int) (int, error) {
	result, err := s.groupRep.DeleteApplicationGroup(ctx, idList)
	if err != nil {
		return 0, errors.WithMessage(err, "delete application group")
	}

	return result, nil
}

func (s ApplicationGroup) convertApplicationGroup(req entity.ApplicationGroup) domain.ApplicationGroup {
	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}
	result := domain.ApplicationGroup{
		Id:          req.Id,
		Name:        req.Name,
		Description: desc,
		DomainId:    req.DomainId,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}

	return result
}
