package controller

import (
	"context"
	"fmt"

	"isp-system-service/domain"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/apierrors"
	"google.golang.org/grpc/codes"
)

//nolint:interfacebloat
type ApplicationService interface {
	GetById(ctx context.Context, appId int) (*domain.Application, error)
	GetByToken(ctx context.Context, token string) (*domain.GetApplicationByTokenResponse, error)
	GetByIdList(ctx context.Context, idList []int) ([]domain.Application, error)
	GetByAppGroup(ctx context.Context, appGroupId int) ([]domain.Application, error)
	CreateUpdate(ctx context.Context, req domain.ApplicationCreateUpdateRequest) (*domain.Application, error)
	Delete(ctx context.Context, idList []int) (int, error)
	NextId(ctx context.Context) (int, error)
	GetAll(ctx context.Context) ([]domain.Application, error)
	Create(ctx context.Context, req domain.CreateApplicationRequest) (*domain.Application, error)
	Update(ctx context.Context, req domain.UpdateApplicationRequest) (*domain.Application, error)
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
//
//	@Tags			application
//	@Summary		Получить приложение по идентификатору
//	@Description	Возвращает описание приложения по его идентификатору
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.Identity	true	"Идентификатор приложения"
//	@Success		200		{object}	domain.Application
//	@Failure		400		{object}	apierrors.Error
//	@Failure		404		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/application/get_application_by_id [POST]
func (c Application) GetById(ctx context.Context, req domain.Identity) (*domain.Application, error) {
	result, err := c.service.GetById(ctx, req.Id)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeApplicationNotFound,
			fmt.Sprintf("application with id %d not found", req.Id),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// GetByToken godoc
//
//	@Tags			application
//	@Summary		Получить идентификатор приложения и группы по токену
//	@Description	Возвращает идентификатор приложения и группы по токену
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.GetApplicationByTokenRequest	true	"Тело запроса"
//	@Success		200		{object}	domain.GetApplicationByTokenResponse
//	@Failure		400		{object}	apierrors.Error
//	@Failure		404		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/application/get_application_by_token [POST]
func (c Application) GetByToken(ctx context.Context, req domain.GetApplicationByTokenRequest) (*domain.GetApplicationByTokenResponse, error) {
	result, err := c.service.GetByToken(ctx, req.Token)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeApplicationNotFound,
			fmt.Sprintf("application with token '%s' not found", req.Token),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// GetByIdList godoc
//
//	@Tags			application
//	@Summary		Получить список приложений
//	@Description	Возвращает массив приложений их идентификаторам
//	@Accept			json
//	@Produce		json
//	@Param			body	body		[]integer	false	"Массив идентификаторов приложений"
//	@Success		200		{array}		domain.Application
//	@Failure		500		{object}	apierrors.Error
//	@Router			/application/get_applications [POST]
func (c Application) GetByIdList(ctx context.Context, req []int) ([]domain.Application, error) {
	return c.service.GetByIdList(ctx, req)
}

// GetByAppGroup godoc
//
//	@Tags			application
//	@Summary		Получить список приложений по идентификатору группы
//	@Description	Возвращает список приложений по идентификатору группы
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.Identity	true	"Тело запроса"
//	@Success		200		{array}		domain.Application
//	@Failure		400		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/application/get_by_app_group [POST]
func (c Application) GetByAppGroup(ctx context.Context, req domain.Identity) ([]domain.Application, error) {
	return c.service.GetByAppGroup(ctx, req.Id)
}

// CreateUpdate godoc
//
//	@Tags			application
//	@Summary		Создать/обновить приложение
//	@Description	Если приложение с такими идентификатором существует, то обновляет данные, если нет, то добавляет данные в базу
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.ApplicationCreateUpdateRequest	true	"Объект приложения"
//	@Success		200		{object}	domain.Application
//	@Failure		400		{object}	apierrors.Error
//	@Failure		404		{object}	apierrors.Error
//	@Failure		409		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/application/create_update_application [POST]
func (c Application) CreateUpdate(ctx context.Context, req domain.ApplicationCreateUpdateRequest) (*domain.Application, error) {
	result, err := c.service.CreateUpdate(ctx, req)
	switch {
	case errors.Is(err, domain.ErrAppGroupNotFound):
		return nil, apierrors.NewBusinessError(
			domain.ErrCodeAppGroupNotFound,
			fmt.Sprintf("service with id %d not found", req.ApplicationGroupId),
			err,
		)
	case errors.Is(err, domain.ErrApplicationDuplicateName):
		return nil, apierrors.New(
			codes.AlreadyExists,
			domain.ErrCodeApplicationDuplicateName,
			fmt.Sprintf("application with name %s already exists", req.Name),
			err,
		)
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeApplicationNotFound,
			fmt.Sprintf("application with id %d not found", req.Id),
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
//	@Tags			application
//	@Summary		Удалить приложения
//	@Description	Удаляет приложения по списку их идентификаторов, возвращает количество удаленных приложений
//	@Accept			json
//	@Produce		json
//	@Param			body	body		[]integer	false	"Массив идентификаторов приложений"
//	@Success		200		{object}	domain.DeleteResponse
//	@Failure		400		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/application/delete_applications [POST]
func (c Application) Delete(ctx context.Context, req []int) (*domain.DeleteResponse, error) {
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

// NextId godoc
//
//	@Tags			application
//	@Summary		Получить следующий идентификатор приложения
//	@Description	Возвращает следующий идентификатор приложения
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	integer
//	@Failure		500	{object}	apierrors.Error
//	@Router			/application/next_id [POST]
func (c Application) NextId(ctx context.Context) (int, error) {
	return c.service.NextId(ctx)
}

// GetAll godoc
//
//	@Tags			application
//	@Summary		Получить список приложений
//	@Description	Возвращает список приложений
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		domain.Application
//	@Failure		500	{object}	apierrors.Error
//	@Router			/application/get_all [POST]
func (c Application) GetAll(ctx context.Context) ([]domain.Application, error) {
	return c.service.GetAll(ctx)
}

// Create godoc
//
//	@Tags			application
//	@Summary		Создать приложение
//	@Description	Если приложение с такими идентификатором или связкой `applicationGroupId`-`name` существует, то возвращает ошибку
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.CreateApplicationRequest	true	"Объект приложения"
//	@Success		200		{object}	domain.Application
//	@Failure		400		{object}	apierrors.Error
//	@Failure		409		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/application/create_application [POST]
func (c Application) Create(ctx context.Context, req domain.CreateApplicationRequest) (*domain.Application, error) {
	result, err := c.service.Create(ctx, req)
	switch {
	case errors.Is(err, domain.ErrAppGroupNotFound):
		return nil, apierrors.NewBusinessError(
			domain.ErrCodeAppGroupNotFound,
			fmt.Sprintf("application group with id %d not found", req.ApplicationGroupId),
			err,
		)
	case errors.Is(err, domain.ErrApplicationDuplicateName):
		return nil, apierrors.New(
			codes.AlreadyExists,
			domain.ErrCodeApplicationDuplicateName,
			fmt.Sprintf("application with name %s already exists", req.Name),
			err,
		)
	case errors.Is(err, domain.ErrApplicationDuplicateId):
		return nil, apierrors.New(
			codes.AlreadyExists,
			domain.ErrCodeApplicationDuplicateId,
			fmt.Sprintf("application with id %d already exists", req.Id),
			err,
		)
	default:
		return result, err
	}
}

// Update godoc
//
//	@Tags			application
//	@Summary		Обновить приложение
//	@Description	Если приложение с связкой `applicationGroupId`-`name` существует или приложение не найдено, то возвращает ошибку
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.UpdateApplicationRequest	true	"Объект приложения"
//	@Success		200		{object}	domain.Application
//	@Failure		400		{object}	apierrors.Error
//	@Failure		404		{object}	apierrors.Error
//	@Failure		409		{object}	apierrors.Error
//	@Failure		500		{object}	apierrors.Error
//	@Router			/application/update_application [POST]
func (c Application) Update(ctx context.Context, req domain.UpdateApplicationRequest) (*domain.Application, error) {
	result, err := c.service.Update(ctx, req)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeApplicationNotFound,
			fmt.Sprintf("application with id %d not found", req.OldId),
			err,
		)
	case errors.Is(err, domain.ErrApplicationDuplicateName):
		return nil, apierrors.New(
			codes.AlreadyExists,
			domain.ErrCodeApplicationDuplicateName,
			fmt.Sprintf("application with name %s already exists", req.Name),
			err,
		)
	case errors.Is(err, domain.ErrApplicationDuplicateId):
		return nil, apierrors.New(
			codes.AlreadyExists,
			domain.ErrCodeApplicationDuplicateId,
			fmt.Sprintf("application with id %d already exists", req.NewId),
			err,
		)
	default:
		return result, err
	}
}
