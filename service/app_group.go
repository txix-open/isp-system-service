package service

import (
	"context"
	"isp-system-service/domain"
	"isp-system-service/entity"

	"github.com/pkg/errors"
)

type AppGroup struct {
	repo AppGroupRepo
}

func NewAppGroup(repo AppGroupRepo) AppGroup {
	return AppGroup{
		repo: repo,
	}
}

func (s AppGroup) Create(ctx context.Context, req domain.CreateAppGroupRequest) (*domain.AppGroup, error) {
	appGroup, err := s.repo.CreateAppGroup(ctx, req.Name, req.Description, 1)
	if err != nil {
		return nil, errors.WithMessage(err, "create appGroup")
	}
	converted := s.convertAppGroup(*appGroup)
	return &converted, nil
}

func (s AppGroup) Update(ctx context.Context, req domain.UpdateAppGroupRequest) (*domain.AppGroup, error) {
	appGroup, err := s.repo.UpdateAppGroup(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		return nil, errors.WithMessage(err, "update appGroup")
	}
	converted := s.convertAppGroup(*appGroup)
	return &converted, nil
}

func (s AppGroup) DeleteList(ctx context.Context, req domain.IdListRequest) (*domain.DeleteResponse, error) {
	deleted, err := s.repo.DeleteAppGroup(ctx, req.IdList)
	if err != nil {
		return nil, errors.WithMessage(err, "delete appGroup")
	}
	return &domain.DeleteResponse{
		Deleted: deleted,
	}, nil
}

func (s AppGroup) GetByIdList(ctx context.Context, idList []int) ([]domain.AppGroup, error) {
	appGroups, err := s.repo.GetAppGroupByIdList(ctx, idList)
	if err != nil {
		return nil, errors.WithMessage(err, "get appGroups by id list")
	}
	result := make([]domain.AppGroup, 0, len(appGroups))
	for _, appGroup := range appGroups {
		result = append(result, s.convertAppGroup(appGroup))
	}
	return result, nil
}

func (s AppGroup) convertAppGroup(appGroup entity.AppGroup) domain.AppGroup {
	return domain.AppGroup{
		Id:          appGroup.Id,
		Name:        appGroup.Name,
		Description: appGroup.Description.String,
		CreatedAt:   appGroup.CreatedAt,
		UpdatedAt:   appGroup.UpdatedAt,
	}
}
