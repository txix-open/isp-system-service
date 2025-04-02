package controller

import (
	"context"
	"isp-system-service/domain"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppGroupService interface {
	Create(ctx context.Context, req domain.CreateAppGroupRequest) (*domain.AppGroup, error)
	Update(ctx context.Context, req domain.UpdateAppGroupRequest) (*domain.AppGroup, error)
	DeleteList(ctx context.Context, req domain.IdListRequest) (*domain.DeleteResponse, error)
	GetByIdList(ctx context.Context, idList []int) ([]domain.AppGroup, error)
}

type AppGroup struct {
	service AppGroupService
}

func NewAppGroup(service AppGroupService) AppGroup {
	return AppGroup{
		service: service,
	}
}

// Create godoc
//
//	@Tags			application_group
//	@Summary		Создать группу приложений
//	@Description	Если группа приложений таким именем существует, возвращает ошибку
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.CreateAppGroupRequest	true	"Объект группы приложений"
//	@Success		200		{object}	domain.AppGroup
//	@Failure		400		{object}	domain.GrpcError
//	@Failure		409		{object}	domain.GrpcError
//	@Failure		500		{object}	domain.GrpcError
//	@Router			/application_group/create [POST]
func (c AppGroup) Create(ctx context.Context, req domain.CreateAppGroupRequest) (*domain.AppGroup, error) {
	result, err := c.service.Create(ctx, req)
	switch {
	case errors.Is(err, domain.ErrAppGroupDuplicateName):
		return nil, status.Errorf(codes.AlreadyExists, "application group with name %s already exists", req.Name)
	default:
		return result, err
	}
}

// Update godoc
//
//	@Tags			application_group
//	@Summary		Обновить группу приложений
//	@Description	Если группа приложений таким именем существует или группы приложений с указанным id не существует, возвращает ошибку
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.UpdateAppGroupRequest	true	"Объект группы приложений"
//	@Success		200		{object}	domain.AppGroup
//	@Failure		400		{object}	domain.GrpcError
//	@Failure		404		{object}	domain.GrpcError
//	@Failure		409		{object}	domain.GrpcError
//	@Failure		500		{object}	domain.GrpcError
//	@Router			/application_group/update [POST]
func (c AppGroup) Update(ctx context.Context, req domain.UpdateAppGroupRequest) (*domain.AppGroup, error) {
	result, err := c.service.Update(ctx, req)
	switch {
	case errors.Is(err, domain.ErrAppGroupDuplicateName):
		return nil, status.Errorf(codes.AlreadyExists, "application group with name %s already exists", req.Name)
	case errors.Is(err, domain.ErrAppGroupNotFound):
		return nil, status.Errorf(codes.NotFound, "application group with id %d not found", req.Id)
	default:
		return result, err
	}
}

// DeleteList godoc
//
//	@Tags			application_group
//	@Summary		Удалить группы приложений
//	@Description	Удаляет группы приложений по списку их идентификаторов, возвращает количество удаленных групп приложений
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.IdListRequest	true	"список идентификаторов групп приложений"
//	@Success		200		{object}	domain.DeleteResponse
//	@Failure		400		{object}	domain.GrpcError
//	@Failure		500		{object}	domain.GrpcError
//	@Router			/application_group/delete_list [POST]
func (c AppGroup) DeleteList(ctx context.Context, req domain.IdListRequest) (*domain.DeleteResponse, error) {
	return c.service.DeleteList(ctx, req)
}

// GetByIdList godoc
//
//	@Tags			application_group
//	@Summary		Получить группы приложений по списку идентификаторов
//	@Description	Возвращает группы приложений с указанными идентификаторами
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.IdListRequest	true	"список идентификаторов групп приложений"
//	@Success		200		{array}		domain.AppGroup
//	@Failure		400		{object}	domain.GrpcError
//	@Failure		500		{object}	domain.GrpcError
//	@Router			/application_group/get_by_id_list [POST]
func (c AppGroup) GetByIdList(ctx context.Context, req domain.IdListRequest) ([]domain.AppGroup, error) {
	return c.service.GetByIdList(ctx, req.IdList)
}
