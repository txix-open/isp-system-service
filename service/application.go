package service

import (
	"context"

	"isp-system-service/domain"
	"isp-system-service/entity"

	"github.com/pkg/errors"
)

type ApplicationDeleteTx interface {
	DeleteApplicationByIdList(ctx context.Context, idList []int) (int, error)
}

type ApplicationTxRunner interface {
	ApplicationDeleteTx(ctx context.Context, tx func(ctx context.Context, tx ApplicationDeleteTx) error) error
}

type Application struct {
	txRunner    ApplicationTxRunner
	appRepo     ApplicationRepo
	serviceRepo AppGroupRepo
	tokenRepo   TokenRepo
}

func NewApplication(
	txRunner ApplicationTxRunner,
	applicationRepo ApplicationRepo,
	appGroupRepo AppGroupRepo,
	tokenRepo TokenRepo,
) Application {
	return Application{
		txRunner:    txRunner,
		appRepo:     applicationRepo,
		serviceRepo: appGroupRepo,
		tokenRepo:   tokenRepo,
	}
}

func (s Application) GetById(ctx context.Context, appId int) (*domain.Application, error) {
	app, err := s.appRepo.GetApplicationById(ctx, appId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	return new(s.convertApplication(*app)), nil
}

func (s Application) GetByToken(ctx context.Context, tokenStr string) (*domain.GetApplicationByTokenResponse, error) {
	token, err := s.tokenRepo.GetTokenById(ctx, tokenStr)
	if err != nil {
		return nil, errors.WithMessage(err, "get token by id")
	}
	if token == nil {
		return nil, domain.ErrApplicationNotFound
	}

	app, err := s.appRepo.GetApplicationById(ctx, token.AppId)
	if err != nil {
		return nil, errors.WithMessage(err, "get app by id")
	}
	return &domain.GetApplicationByTokenResponse{
		ApplicationId:      app.Id,
		ApplicationGroupId: app.ApplicationGroupId,
	}, nil
}

func (s Application) GetByIdList(ctx context.Context, idList []int) ([]domain.Application, error) {
	apps, err := s.appRepo.GetApplicationByIdList(ctx, idList)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id list")
	}

	result := make([]domain.Application, len(apps))
	for i, a := range apps {
		result[i] = s.convertApplication(a)
	}
	return result, nil
}

func (s Application) GetByAppGroup(ctx context.Context, appGroupId int) ([]domain.Application, error) {
	apps, err := s.appRepo.GetApplicationByAppGroupIdList(ctx, []int{appGroupId})
	if err != nil {
		return nil, errors.WithMessage(err, "get apps by app group ids")
	}

	result := make([]domain.Application, len(apps))
	for i, a := range apps {
		result[i] = s.convertApplication(a)
	}
	return result, nil
}

func (s Application) CreateUpdate(ctx context.Context, req domain.ApplicationCreateUpdateRequest) (*domain.Application, error) {
	if req.Id != 0 {
		app, err := s.appRepo.UpdateApplication(ctx, req.Id, req.Name, req.Description)
		if err != nil {
			return nil, errors.WithMessage(err, "update application")
		}

		return new(s.convertApplication(*app)), nil
	}

	appId, err := s.appRepo.NextApplicationId(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "next app id")
	}

	app, err := s.appRepo.CreateApplication(ctx, appId, req.Name, req.Description, req.ApplicationGroupId, req.Type)
	if err != nil {
		return nil, errors.WithMessage(err, "create application")
	}

	return new(s.convertApplication(*app)), nil
}

func (s Application) Delete(ctx context.Context, idList []int) (int, error) {
	count := 0
	err := s.txRunner.ApplicationDeleteTx(ctx, func(ctx context.Context, tx ApplicationDeleteTx) error {
		deletedApp, err := tx.DeleteApplicationByIdList(ctx, idList)
		if err != nil {
			return errors.WithMessage(err, "delete application by id list")
		}

		count = deletedApp
		return nil
	})
	if err != nil {
		return 0, errors.WithMessage(err, "transaction application delete")
	}

	return count, nil
}

func (s Application) EnrichWithTokens(ctx context.Context, apps []entity.Application) ([]*domain.ApplicationWithTokens, error) {
	if len(apps) == 0 {
		return []*domain.ApplicationWithTokens{}, nil
	}

	appIdList := make([]int, len(apps))
	resultByAppId := make(map[int]*domain.ApplicationWithTokens, len(apps))
	result := make([]*domain.ApplicationWithTokens, len(apps))
	for i, a := range apps {
		appIdList[i] = a.Id
		awt := &domain.ApplicationWithTokens{
			App:    s.convertApplication(a),
			Tokens: make([]domain.Token, 0),
		}
		resultByAppId[a.Id] = awt
		result[i] = awt
	}

	tokens, err := s.tokenRepo.GetTokenByAppIdList(ctx, appIdList)
	if err != nil {
		return nil, errors.WithMessage(err, "get token by app_id list")
	}

	for _, token := range tokens {
		r := resultByAppId[token.AppId]
		r.Tokens = append(r.Tokens, domain.Token(token))
	}

	return result, nil
}

func (s Application) NextId(ctx context.Context) (int, error) {
	nextId, err := s.appRepo.NextApplicationId(ctx)
	if err != nil {
		return 0, errors.WithMessage(err, "get next application id")
	}
	return nextId, nil
}

func (s Application) GetAll(ctx context.Context) ([]domain.Application, error) {
	apps, err := s.appRepo.GetAllApplications(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "get applications list")
	}

	result := make([]domain.Application, 0, len(apps))
	for _, app := range apps {
		result = append(result, s.convertApplication(app))
	}
	return result, nil
}

func (s Application) Create(ctx context.Context, req domain.CreateApplicationRequest) (*domain.Application, error) {
	app, err := s.appRepo.CreateApplication(ctx, req.Id, req.Name, req.Description, req.ApplicationGroupId, req.Type)
	if err != nil {
		return nil, errors.WithMessage(err, "create application")
	}

	return new(s.convertApplication(*app)), nil
}

func (s Application) Update(ctx context.Context, req domain.UpdateApplicationRequest) (*domain.Application, error) {
	app, err := s.appRepo.UpdateApplicationWithNewId(ctx, req.OldId, req.NewId, req.Name, req.Description)
	if err != nil {
		return nil, errors.WithMessage(err, "update application")
	}

	return new(s.convertApplication(*app)), nil
}

func (s Application) convertApplication(req entity.Application) domain.Application {
	return domain.Application{
		Id:                 req.Id,
		Name:               req.Name,
		Description:        req.Description.String,
		ApplicationGroupId: req.ApplicationGroupId,
		Type:               req.Type,
		CreatedAt:          req.CreatedAt,
		UpdatedAt:          req.UpdatedAt,
	}
}
