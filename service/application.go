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
	domainRepo  DomainRepo
	serviceRepo AppGroupRepo
	tokenRepo   TokenRepo
}

func NewApplication(
	txRunner ApplicationTxRunner,
	applicationRepo ApplicationRepo,
	domainRepo DomainRepo,
	appGroupRepo AppGroupRepo,
	tokenRepo TokenRepo,
) Application {
	return Application{
		txRunner:    txRunner,
		appRepo:     applicationRepo,
		domainRepo:  domainRepo,
		serviceRepo: appGroupRepo,
		tokenRepo:   tokenRepo,
	}
}

func (s Application) GetById(ctx context.Context, appId int) (*domain.ApplicationWithTokens, error) {
	application, err := s.appRepo.GetApplicationById(ctx, appId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	arr, err := s.EnrichWithTokens(ctx, []entity.Application{*application})
	if err != nil {
		return nil, errors.WithMessage(err, "enrich application with tokens")
	}

	return arr[0], nil
}

func (s Application) GetByIdList(ctx context.Context, idList []int) ([]*domain.ApplicationWithTokens, error) {
	res, err := s.appRepo.GetApplicationByIdList(ctx, idList)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id list")
	}

	return s.EnrichWithTokens(ctx, res)
}

func (s Application) GetByServiceId(ctx context.Context, id int) ([]*domain.ApplicationWithTokens, error) {
	arr, err := s.appRepo.GetApplicationByAppGroupIdList(ctx, []int{id})
	if err != nil {
		return nil, errors.WithMessage(err, "get application by app_group_id")
	}

	return s.EnrichWithTokens(ctx, arr)
}

func (s Application) SystemTree(ctx context.Context, systemId int) ([]*domain.DomainWithService, error) {
	domainEntityList, err := s.domainRepo.GetDomainBySystemId(ctx, systemId)
	if err != nil {
		return nil, errors.WithMessage(err, "get domains by system_id")
	}
	if len(domainEntityList) == 0 {
		return []*domain.DomainWithService{}, nil
	}

	result := make([]*domain.DomainWithService, len(domainEntityList))
	resultByDomainId := make(map[int]*domain.DomainWithService, len(domainEntityList))
	domainIdList := make([]int, len(domainEntityList))
	for i, domainEntity := range domainEntityList {
		domainIdList[i] = domainEntity.Id
		description := ""
		if domainEntity.Description != nil {
			description = *domainEntity.Description
		}
		r := &domain.DomainWithService{
			Id:          domainEntity.Id,
			Name:        domainEntity.Name,
			Description: description,
			Services:    make([]*domain.ServiceWithApps, 0),
		}
		resultByDomainId[domainEntity.Id] = r
		result[i] = r
	}

	serviceEntityList, err := s.serviceRepo.GetAppGroupByDomainId(ctx, domainIdList)
	if err != nil {
		return nil, errors.WithMessage(err, "get service by domain_id")
	}

	serviceIdList := make([]int, len(serviceEntityList))
	resultServiceByServiceId := make(map[int]*domain.ServiceWithApps, len(serviceEntityList))
	for i, serviceEntity := range serviceEntityList {
		serviceIdList[i] = serviceEntity.Id
		description := ""
		if serviceEntity.Description != nil {
			description = *serviceEntity.Description
		}
		resultService := &domain.ServiceWithApps{
			Id:          serviceEntity.Id,
			Name:        serviceEntity.Name,
			Description: description,
			Apps:        make([]*domain.ApplicationSimple, 0),
		}
		r := resultByDomainId[serviceEntity.DomainId]
		r.Services = append(r.Services, resultService)
		resultServiceByServiceId[serviceEntity.Id] = resultService
	}

	applicationEntityList, err := s.appRepo.GetApplicationByAppGroupIdList(ctx, serviceIdList)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by service_id")
	}

	for _, applicationEntity := range applicationEntityList {
		resultService := resultServiceByServiceId[applicationEntity.ApplicationGroupId]
		description := ""
		if applicationEntity.Description != nil {
			description = *applicationEntity.Description
		}
		resultService.Apps = append(resultService.Apps, &domain.ApplicationSimple{
			Id:          applicationEntity.Id,
			Name:        applicationEntity.Name,
			Type:        applicationEntity.Type,
			Description: description,
			Tokens:      make([]domain.Token, 0),
		})
	}

	return result, nil
}

func (s Application) CreateUpdate(ctx context.Context, req domain.ApplicationCreateUpdateRequest) (*domain.ApplicationWithTokens, error) {
	existed, err := s.appRepo.GetApplicationByNameAndAppGroupId(ctx, req.Name, req.ServiceId)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
	case err != nil:
		return nil, errors.WithMessage(err, "get application by name and service_id")
	}

	_, err = s.serviceRepo.GetAppGroupById(ctx, req.ServiceId)
	if err != nil {
		return nil, errors.WithMessage(err, "get service by id")
	}

	if existed != nil && (existed.Id != req.Id || req.Id == 0) {
		return nil, domain.ErrApplicationDuplicateName
	}

	if req.Id == 0 {
		_, err = s.appRepo.GetApplicationByNameAndAppGroupId(ctx, req.Name, req.ServiceId)
		if err != nil {
			return nil, errors.WithMessage(err, "get application by name and service_id")
		}

		appId, err := s.appRepo.NextApplicationId(ctx)
		if err != nil {
			return nil, errors.WithMessage(err, "next app id")
		}

		app, err := s.appRepo.CreateApplication(ctx, appId, req.Name, req.Description, req.ServiceId, req.Type)
		if err != nil {
			return nil, errors.WithMessage(err, "create application")
		}

		result, err := s.EnrichWithTokens(ctx, []entity.Application{*app})
		if err != nil {
			return nil, errors.WithMessage(err, "enrich application with tokens")
		}

		return result[0], nil
	}

	_, err = s.appRepo.GetApplicationById(ctx, req.Id)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	app, err := s.appRepo.UpdateApplication(ctx, req.Id, req.Name, req.Description)
	if err != nil {
		return nil, errors.WithMessage(err, "update application")
	}

	result, err := s.EnrichWithTokens(ctx, []entity.Application{*app})
	if err != nil {
		return nil, errors.WithMessage(err, "enrich application with tokens")
	}

	return result[0], nil
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

func (s Application) convertApplication(req entity.Application) domain.Application {
	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}

	return domain.Application{
		Id:          req.Id,
		Name:        req.Name,
		Description: desc,
		ServiceId:   req.ApplicationGroupId,
		Type:        req.Type,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}
}
