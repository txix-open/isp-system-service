package service

import (
	"context"

	"isp-system-service/domain"
	"isp-system-service/entity"

	"github.com/pkg/errors"
)

type ApplicationTokenCreator interface {
	CreateApplicationToken() (string, error)
}

type TokenAppEnricher interface {
	EnrichWithTokens(ctx context.Context, apps []entity.Application) ([]*domain.ApplicationWithTokens, error)
}

type TokenCreateTx interface {
	SaveToken(ctx context.Context, token string, appId int, expireTime int) (*entity.Token, error)
}

type TokenRevokeTx interface {
	DeleteToken(ctx context.Context, tokens []string) (int, error)
}

type TokenTxRunner interface {
	TokenCreateTx(ctx context.Context, tx func(ctx context.Context, tx TokenCreateTx) error) error
	TokenRevokeTx(ctx context.Context, tx func(ctx context.Context, tx TokenRevokeTx) error) error
}

type Token struct {
	jwt         ApplicationTokenCreator
	appEnrich   TokenAppEnricher
	tx          TokenTxRunner
	appRepo     ApplicationRepo
	domainRepo  DomainRepo
	serviceRepo ServiceRepo
	tokenRepo   TokenRepo
}

func NewToken(
	jwtGenerate ApplicationTokenCreator,
	appEnrich TokenAppEnricher,
	tx TokenTxRunner,
	appRepo ApplicationRepo,
	domainRepo DomainRepo,
	serviceRepo ServiceRepo,
	tokenRepo TokenRepo,
) Token {
	return Token{
		appEnrich:   appEnrich,
		jwt:         jwtGenerate,
		tx:          tx,
		appRepo:     appRepo,
		domainRepo:  domainRepo,
		serviceRepo: serviceRepo,
		tokenRepo:   tokenRepo,
	}
}

func (s Token) GetByAppId(ctx context.Context, appId int) ([]domain.Token, error) {
	tokenEntity, err := s.tokenRepo.GetTokenByAppIdList(ctx, []int{appId})
	if err != nil {
		return nil, errors.WithMessage(err, "get token by app_id list")
	}

	result := make([]domain.Token, len(tokenEntity))
	for i, token := range tokenEntity {
		result[i] = domain.Token(token)
	}

	return result, nil
}

func (s Token) Create(ctx context.Context, req domain.TokenCreateRequest) (*domain.ApplicationWithTokens, error) {
	applicationEntity, err := s.appRepo.GetApplicationById(ctx, req.AppId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	serviceEntity, err := s.serviceRepo.GetServiceById(ctx, applicationEntity.ServiceId)
	if err != nil {
		return nil, errors.WithMessage(err, "get service by id")
	}

	_, err = s.domainRepo.GetDomainById(ctx, serviceEntity.DomainId)
	if err != nil {
		return nil, errors.WithMessage(err, "get domain by id")
	}

	token, err := s.jwt.CreateApplicationToken()
	if err != nil {
		return nil, errors.WithMessage(err, "create application token")
	}

	err = s.tx.TokenCreateTx(ctx, func(ctx context.Context, tx TokenCreateTx) error {
		_, err = tx.SaveToken(ctx, token, req.AppId, req.ExpireTimeMs)
		if err != nil {
			return errors.WithMessage(err, "tx save token")
		}

		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "token create transaction")
	}

	arr, err := s.appEnrich.EnrichWithTokens(ctx, []entity.Application{*applicationEntity})
	if err != nil {
		return nil, errors.WithMessage(err, "application enrich with tokens")
	}

	return arr[0], nil
}

func (s Token) Revoke(ctx context.Context, req domain.TokenRevokeRequest) (*domain.ApplicationWithTokens, error) {
	app, err := s.appRepo.GetApplicationById(ctx, req.AppId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	_, err = s.revokeTokens(ctx, req.Tokens)
	if err != nil {
		return nil, errors.WithMessage(err, "reboke tokens")
	}

	res, err := s.appEnrich.EnrichWithTokens(ctx, []entity.Application{*app})
	if err != nil {
		return nil, errors.WithMessage(err, "application enrich")
	}

	return res[0], nil
}

func (s Token) RevokeByAppId(ctx context.Context, appId int) (*domain.DeleteResponse, error) {
	tokens, err := s.tokenRepo.GetTokenByAppIdList(ctx, []int{appId})
	if err != nil {
		return nil, errors.WithMessage(err, "get token by app_id list")
	}

	tokenIdList := make([]string, len(tokens))
	for i, t := range tokens {
		tokenIdList[i] = t.Token
	}

	return s.revokeTokens(ctx, tokenIdList)
}

func (s Token) revokeTokens(ctx context.Context, tokens []string) (*domain.DeleteResponse, error) {
	if len(tokens) == 0 {
		return &domain.DeleteResponse{Deleted: 0}, nil
	}

	var count int
	err := s.tx.TokenRevokeTx(ctx, func(ctx context.Context, tx TokenRevokeTx) error {
		deleted, err := tx.DeleteToken(ctx, tokens)
		if err != nil {
			return errors.WithMessage(err, "tx delete token")
		}

		count = deleted
		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "token revoke transaction")
	}

	return &domain.DeleteResponse{
		Deleted: count,
	}, nil
}
