package service

import (
	"context"

	"github.com/pkg/errors"
	"isp-system-service/domain"
	"isp-system-service/entity"
)

type IApplicationApplicationRep interface {
	GetApplicationById(ctx context.Context, id int) (*entity.Application, error)
	GetApplicationByIdList(ctx context.Context, idList []int) ([]entity.Application, error)
	GetApplicationByApplicationGroupIdList(ctx context.Context, serviceIdList []int) ([]entity.Application, error)
	GetApplicationByNameAndApplicationGroupId(ctx context.Context, name string, serviceId int) (*entity.Application, error)
	CreateApplication(ctx context.Context, name string, desc string, serviceId int, appType string) (*entity.Application, error)
	UpdateApplication(ctx context.Context, id int, name string, description string) (*entity.Application, error)
}

type IApplicationTokenRep interface {
	GetTokenByAppIdList(ctx context.Context, appIdList []int) ([]entity.Token, error)
}

type IApplicationServiceRep interface {
	GetApplicationGroupById(ctx context.Context, id int) (*entity.ApplicationGroup, error)
	GetApplicationGroupByDomainId(ctx context.Context, domainIdList []int) ([]entity.ApplicationGroup, error)
}

type IApplicationDomainRep interface {
	GetDomainBySystemId(ctx context.Context, systemId int) ([]entity.Domain, error)
}

type IApplicationDeleteTx interface {
	DeleteApplicationByIdList(ctx context.Context, idList []int) (int, error)
}

type IApplicationTxRunner interface {
	ApplicationDeleteTx(ctx context.Context, tx func(ctx context.Context, tx IApplicationDeleteTx) error) error
}

type Application struct {
	tx             IApplicationTxRunner
	applicationRep IApplicationApplicationRep
	domainRep      IApplicationDomainRep
	groupRep       IApplicationServiceRep
	tokenRep       IApplicationTokenRep
}

func NewApplication(
	tx IApplicationTxRunner,
	applicationRep IApplicationApplicationRep,
	domainRep IApplicationDomainRep,
	groupRep IApplicationServiceRep,
	tokenRep IApplicationTokenRep,
) Application {
	return Application{
		tx:             tx,
		applicationRep: applicationRep,
		domainRep:      domainRep,
		groupRep:       groupRep,
		tokenRep:       tokenRep,
	}
}

func (s Application) GetById(ctx context.Context, appId int) (*domain.ApplicationWithTokens, error) {
	application, err := s.applicationRep.GetApplicationById(ctx, appId)
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
	res, err := s.applicationRep.GetApplicationByIdList(ctx, idList)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id list")
	}

	return s.EnrichWithTokens(ctx, res)
}

func (s Application) GetByApplicationGroupId(ctx context.Context, id int) ([]*domain.ApplicationWithTokens, error) {
	arr, err := s.applicationRep.GetApplicationByApplicationGroupIdList(ctx, []int{id})
	if err != nil {
		return nil, errors.WithMessage(err, "get application by application_group_id")
	}

	return s.EnrichWithTokens(ctx, arr)
}

func (s Application) SystemTree(ctx context.Context, systemId int) ([]*domain.DomainWithApplicationGroup, error) {
	domainEntityList, err := s.domainRep.GetDomainBySystemId(ctx, systemId)
	if err != nil {
		return nil, errors.WithMessage(err, "get domains by system_id")
	}
	if len(domainEntityList) == 0 {
		return []*domain.DomainWithApplicationGroup{}, nil
	}

	result := make([]*domain.DomainWithApplicationGroup, len(domainEntityList))
	resultByDomainId := make(map[int]*domain.DomainWithApplicationGroup, len(domainEntityList))
	domainIdList := make([]int, len(domainEntityList))
	for i, domainEntity := range domainEntityList {
		domainIdList[i] = domainEntity.Id
		description := ""
		if domainEntity.Description != nil {
			description = *domainEntity.Description
		}
		r := &domain.DomainWithApplicationGroup{
			Id:               domainEntity.Id,
			Name:             domainEntity.Name,
			Description:      description,
			ApplicationGroup: make([]*domain.ApplicationGroupWithApps, 0),
		}
		resultByDomainId[domainEntity.Id] = r
		result[i] = r
	}

	applicationGroupEntityList, err := s.groupRep.GetApplicationGroupByDomainId(ctx, domainIdList)
	if err != nil {
		return nil, errors.WithMessage(err, "get application group by domain_id")
	}

	applicationGroupIdList := make([]int, len(applicationGroupEntityList))
	resultApplicationGroupById := make(map[int]*domain.ApplicationGroupWithApps, len(applicationGroupEntityList))
	for i, applicationGroupEntity := range applicationGroupEntityList {
		applicationGroupIdList[i] = applicationGroupEntity.Id
		description := ""
		if applicationGroupEntity.Description != nil {
			description = *applicationGroupEntity.Description
		}
		resultApplicationGroup := &domain.ApplicationGroupWithApps{
			Id:          applicationGroupEntity.Id,
			Name:        applicationGroupEntity.Name,
			Description: description,
			Apps:        make([]*domain.ApplicationSimple, 0),
		}
		r := resultByDomainId[applicationGroupEntity.DomainId]
		r.ApplicationGroup = append(r.ApplicationGroup, resultApplicationGroup)
		resultApplicationGroupById[applicationGroupEntity.Id] = resultApplicationGroup
	}

	applicationEntityList, err := s.applicationRep.GetApplicationByApplicationGroupIdList(ctx, applicationGroupIdList)
	if err != nil {
		return nil, errors.WithMessagef(err, "get application by application_group_id")
	}

	for _, applicationEntity := range applicationEntityList {
		resultApplicationGroup := resultApplicationGroupById[applicationEntity.ApplicationGroupId]
		description := ""
		if applicationEntity.Description != nil {
			description = *applicationEntity.Description
		}
		resultApplicationGroup.Apps = append(resultApplicationGroup.Apps, &domain.ApplicationSimple{
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
	existed, err := s.applicationRep.GetApplicationByNameAndApplicationGroupId(ctx, req.Name, req.ApplicationGroupId)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
	case err != nil:
		return nil, errors.WithMessage(err, "get application by name and application_group_id")
	}

	_, err = s.groupRep.GetApplicationGroupById(ctx, req.ApplicationGroupId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application group by id")
	}

	if req.Id == 0 {
		if existed != nil {
			return nil, domain.ErrApplicationDuplicateName
		}

		_, err = s.applicationRep.GetApplicationByNameAndApplicationGroupId(ctx, req.Name, req.ApplicationGroupId)
		if err != nil {
			return nil, errors.WithMessage(err, "get application by name and application_group_id")
		}

		app, err := s.applicationRep.CreateApplication(ctx, req.Name, req.Description, req.ApplicationGroupId, req.Type)
		if err != nil {
			return nil, errors.WithMessage(err, "create application")
		}

		result, err := s.EnrichWithTokens(ctx, []entity.Application{*app})
		if err != nil {
			return nil, errors.WithMessage(err, "enrich application with tokens")
		}

		return result[0], nil
	}

	if existed != nil && existed.Id != req.Id {
		return nil, domain.ErrApplicationDuplicateName
	}

	_, err = s.applicationRep.GetApplicationById(ctx, req.Id)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	app, err := s.applicationRep.UpdateApplication(ctx, req.Id, req.Name, req.Description)
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
	err := s.tx.ApplicationDeleteTx(ctx, func(ctx context.Context, tx IApplicationDeleteTx) error {
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

	tokens, err := s.tokenRep.GetTokenByAppIdList(ctx, appIdList)
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
		Id:                 req.Id,
		Name:               req.Name,
		Description:        desc,
		ApplicationGroupId: req.ApplicationGroupId,
		Type:               req.Type,
		CreatedAt:          req.CreatedAt,
		UpdatedAt:          req.UpdatedAt,
	}
}
