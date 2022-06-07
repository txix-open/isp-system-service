package controller

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"isp-system-service/domain"
)

type ApplicationService interface {
	GetById(ctx context.Context, appId int) (*domain.ApplicationWithTokens, error)
	GetByIdList(ctx context.Context, idList []int) ([]*domain.ApplicationWithTokens, error)
	GetByServiceId(ctx context.Context, id int) ([]*domain.ApplicationWithTokens, error)
	SystemTree(ctx context.Context, systemId int) ([]*domain.DomainWithService, error)
	CreateUpdate(ctx context.Context, req domain.ApplicationCreateUpdateRequest) (*domain.ApplicationWithTokens, error)
	Delete(ctx context.Context, idList []int) (int, error)
}

type Application struct {
	service ApplicationService
}

func NewApplication(service ApplicationService) Application {
	return Application{
		service: service,
	}
}

// GetById godoc
// @Tags application
// @Summary Получить приложение по идентификатору
// @Description  Возвращает описание приложения по его идентификатору
// @Accept  json
// @Produce  json
// @Param body body domain.Identity true "Идентификатор приложения"
// @Success 200 {object} domain.ApplicationWithTokens
// @Failure 400 {object} structure.GrpcError
// @Failure 404 {object} structure.GrpcError
// @Failure 500 {object} structure.GrpcError
// @Router /application/get_application_by_id [POST]
func (c Application) GetById(ctx context.Context, req domain.Identity) (*domain.ApplicationWithTokens, error) {
	result, err := c.service.GetById(ctx, req.Id)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, status.Errorf(codes.NotFound, "application with id %d not found", req.Id)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// GetByIdList godoc
// @Tags application
// @Summary Получить список приложений
// @Description Возвращает массив приложений с токенами по их идентификаторам
// @Accept  json
// @Produce  json
// @Param body body []integer false "Массив идентификаторов приложений"
// @Success 200 {array} domain.ApplicationWithTokens
// @Failure 500 {object} structure.GrpcError
// @Router /application/get_applications [POST]
func (c Application) GetByIdList(ctx context.Context, req []int) ([]*domain.ApplicationWithTokens, error) {
	return c.service.GetByIdList(ctx, req)
}

// GetByServiceId godoc
// @Tags application
// @Summary Получить список приложений по идентификатору сервиса
// @Description Возвращает список приложений по запрошенному идентификатору сервиса
// @Accept  json
// @Produce  json
// @Param body body domain.Identity true "Идентификатор сервиса"
// @Success 200 {array} domain.ApplicationWithTokens
// @Failure 500 {object} structure.GrpcError
// @Router /application/get_applications_by_service_id [POST]
func (c Application) GetByServiceId(ctx context.Context, req domain.Identity) ([]*domain.ApplicationWithTokens, error) {
	return c.service.GetByServiceId(ctx, req.Id)
}

// GetSystemTree godoc
// @Tags application
// @Summary Метод получения системного дерева
// @Description Возвращает описание взаимосвязей сервисов и приложений
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.DomainWithService
// @Failure 500 {object} structure.GrpcError
// @Router /application/get_system_tree [POST]
func (c Application) GetSystemTree(ctx context.Context) ([]*domain.DomainWithService, error) {
	return c.service.SystemTree(ctx, domain.DefaultSystemId)
}

// CreateUpdate godoc
// @Tags application
// @Summary Создать/обновить приложение
// @Description Если приложение с такими идентификатором существует, то обновляет данные, если нет, то добавляет данные в базу
// @Accept  json
// @Produce  json
// @Param body body domain.ApplicationCreateUpdateRequest true "Объект приложения"
// @Success 200 {object} domain.ApplicationWithTokens
// @Failure 400 {object} structure.GrpcError
// @Failure 404 {object} structure.GrpcError
// @Failure 409 {object} structure.GrpcError
// @Failure 500 {object} structure.GrpcError
// @Router /application/create_update_application [POST]
func (c Application) CreateUpdate(ctx context.Context, req domain.ApplicationCreateUpdateRequest) (*domain.ApplicationWithTokens, error) {
	result, err := c.service.CreateUpdate(ctx, req)
	switch {
	case errors.Is(err, domain.ErrServiceNotFound):
		return nil, status.Errorf(codes.InvalidArgument, "service with id %d not found", req.ServiceId)
	case errors.Is(err, domain.ErrApplicationDuplicateName):
		return nil, status.Errorf(codes.AlreadyExists, "application with name %s already exists", req.Name)
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, status.Errorf(codes.NotFound, "Application with id %d not found", req.Id)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// Delete godoc
// @Tags application
// @Summary Удалить приложения
// @Description Удаляет приложения по списку их идентификаторов, возвращает количество удаленных приложений
// @Accept  json
// @Produce  json
// @Param body body []integer false "Массив идентификаторов приложений"
// @Success 200 {object} domain.DeleteResponse
// @Failure 400 {object} structure.GrpcError
// @Failure 500 {object} structure.GrpcError
// @Router /application/delete_applications [POST]
func (c Application) Delete(ctx context.Context, req []int) (*domain.DeleteResponse, error) {
	if len(req) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "At least one id are required")
	}

	result, err := c.service.Delete(ctx, req)
	if err != nil {
		return nil, err
	}

	return &domain.DeleteResponse{
		Deleted: result,
	}, nil
}
