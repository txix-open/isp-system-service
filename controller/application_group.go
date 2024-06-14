package controller

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"isp-system-service/domain"
)

type ApplicationGroupService interface {
	GetById(ctx context.Context, id int) (*domain.ApplicationGroup, error)
	GetByIdList(ctx context.Context, idList []int) ([]domain.ApplicationGroup, error)
	CreateUpdate(ctx context.Context, req domain.ApplicationGroupCreateUpdateRequest) (*domain.ApplicationGroup, error)
	Delete(ctx context.Context, idList []int) (int, error)
}

type ApplicationGroup struct {
	service ApplicationGroupService
}

func NewApplicationGroup(service ApplicationGroupService) ApplicationGroup {
	return ApplicationGroup{
		service: service,
	}
}

// GetById godoc
// @Tags application_group
// @Summary Получить группу приложений по идентификатору
// @Description Возвращает описание группы приложений по ее идентификатору
// @Accept  json
// @Produce  json
// @Param body body domain.Identity true "Идентификатор группы приложений"
// @Success 200 {object} domain.ApplicationGroup
// @Failure 404 {object} domain.GrpcError
// @Failure 500 {object} domain.GrpcError
// @Router /application_group/get_group_by_id [POST]
func (c ApplicationGroup) GetById(ctx context.Context, req domain.Identity) (*domain.ApplicationGroup, error) {
	result, err := c.service.GetById(ctx, req.Id)
	switch {
	case errors.Is(err, domain.ErrApplicationGroupNotFound):
		return nil, status.Errorf(codes.NotFound, "application group with id %d not found", req.Id)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// Get godoc
// @Tags application_group
// @Summary Получить список группы приложений
// @Description Возвращает список групп приложений по их идентификаторам
// @Accept  json
// @Produce  json
// @Param body body []integer false "Массив идентификаторов групп приложений"
// @Success 200 {array} domain.ApplicationGroup
// @Failure 500 {object} domain.GrpcError
// @Router /application_group/get_group [POST]
func (c ApplicationGroup) Get(ctx context.Context, req []int) ([]domain.ApplicationGroup, error) {
	return c.service.GetByIdList(ctx, req)
}

// CreateUpdate godoc
// @Tags application_group
// @Summary Создать/обновить группу приложений
// @Description Если группы приложений с такими идентификатором существует, то обновляет данные, если нет, то добавляет данные в базу
// @Accept  json
// @Produce  json
// @Param body body domain.ApplicationGroupCreateUpdateRequest true "Объект группы приложений"
// @Success 200 {object} domain.ApplicationGroup
// @Failure 400 {object} domain.GrpcError
// @Failure 404 {object} domain.GrpcError
// @Failure 409 {object} domain.GrpcError
// @Failure 500 {object} domain.GrpcError
// @Router /application_group/create_update_group [POST]
func (c ApplicationGroup) CreateUpdate(ctx context.Context, req domain.ApplicationGroupCreateUpdateRequest) (*domain.ApplicationGroup, error) {
	req.DomainId = 1 // Хардкодом проставляем domainId = 1

	result, err := c.service.CreateUpdate(ctx, req)
	switch {
	case errors.Is(err, domain.ErrDomainNotFound):
		return nil, status.Errorf(codes.InvalidArgument, "domain with id %d not found", req.DomainId)
	case errors.Is(err, domain.ErrApplicationGroupNotFound):
		return nil, status.Errorf(codes.NotFound, "application group with id %d not found", req.Id)
	case errors.Is(err, domain.ErrApplicationGroupDuplicateName):
		return nil, status.Errorf(codes.AlreadyExists, "application group with name %s already exists", req.Name)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// Delete godoc
// @Tags application_group
// @Summary Удалить группы приложений
// @Description Удаляет группы приложений по списку их идентификаторов, возвращает количество удаленных групп
// @Accept  json
// @Produce  json
// @Param body body []integer true "Массив идентификаторов групп приложений"
// @Success 200 {object} domain.DeleteResponse
// @Failure 400 {object} domain.GrpcError
// @Failure 500 {object} domain.GrpcError
// @Router /application_group/delete_group [POST]
func (c ApplicationGroup) Delete(ctx context.Context, req []int) (*domain.DeleteResponse, error) {
	if len(req) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one id are required")
	}

	result, err := c.service.Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	return &domain.DeleteResponse{
		Deleted: result,
	}, nil
}
