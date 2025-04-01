package service

import (
	"context"

	"github.com/pkg/errors"
	"isp-system-service/domain"
	"isp-system-service/entity"
)

type ITokenSource interface {
	CreateApplicationToken() (string, error)
}

type ITokenAppEnrich interface {
	EnrichWithTokens(ctx context.Context, apps []entity.Application) ([]*domain.ApplicationWithTokens, error)
}

type ITokenApplicationRep interface { // nolint:iface
	GetApplicationById(ctx context.Context, id int) (*entity.Application, error)
}

type ITokenServiceRep interface {
	GetServiceById(ctx context.Context, id int) (*entity.Service, error)
}

type ITokenDomainRep interface { // nolint:iface
	GetDomainById(ctx context.Context, id int) (*entity.Domain, error)
}

type ITokenTokenRep interface { // nolint:iface
	GetTokenByAppIdList(ctx context.Context, appIdList []int) ([]entity.Token, error)
}

type ITokenCreateTx interface {
	SaveToken(ctx context.Context, token string, appId int, expireTime int) (*entity.Token, error)
}

type ITokenRevokeTx interface {
	DeleteToken(ctx context.Context, tokens []string) (int, error)
}

type ITokenTxRunner interface {
	TokenCreateTx(ctx context.Context, tx func(ctx context.Context, tx ITokenCreateTx) error) error
	TokenRevokeTx(ctx context.Context, tx func(ctx context.Context, tx ITokenRevokeTx) error) error
}

type Token struct {
	jwt        ITokenSource
	appEnrich  ITokenAppEnrich
	tx         ITokenTxRunner
	appRep     ITokenApplicationRep
	domainRep  ITokenDomainRep
	serviceRep ITokenServiceRep
	tokenRep   ITokenTokenRep
}

func NewToken(
	jwtGenerate ITokenSource,
	appEnrich ITokenAppEnrich,
	tx ITokenTxRunner,
	appRep ITokenApplicationRep,
	domainRep ITokenDomainRep,
	serviceRep ITokenServiceRep,
	tokenRep ITokenTokenRep,
) Token {
	return Token{
		appEnrich:  appEnrich,
		jwt:        jwtGenerate,
		tx:         tx,
		appRep:     appRep,
		domainRep:  domainRep,
		serviceRep: serviceRep,
		tokenRep:   tokenRep,
	}
}

func (s Token) GetByAppId(ctx context.Context, appId int) ([]domain.Token, error) {
	tokenEntity, err := s.tokenRep.GetTokenByAppIdList(ctx, []int{appId})
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
	applicationEntity, err := s.appRep.GetApplicationById(ctx, req.AppId)
	if err != nil {
		return nil, errors.WithMessagef(err, "get application by id")
	}

	serviceEntity, err := s.serviceRep.GetServiceById(ctx, applicationEntity.ServiceId)
	if err != nil {
		return nil, errors.WithMessagef(err, "get service by id")
	}

	_, err = s.domainRep.GetDomainById(ctx, serviceEntity.DomainId)
	if err != nil {
		return nil, errors.WithMessagef(err, "get domain by id")
	}

	token, err := s.jwt.CreateApplicationToken()
	if err != nil {
		return nil, errors.WithMessagef(err, "create application token")
	}

	err = s.tx.TokenCreateTx(ctx, func(ctx context.Context, tx ITokenCreateTx) error {
		_, err = tx.SaveToken(ctx, token, req.AppId, req.ExpireTimeMs)
		if err != nil {
			return errors.WithMessagef(err, "tx save token")
		}

		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "token create transaction")
	}

	arr, err := s.appEnrich.EnrichWithTokens(ctx, []entity.Application{*applicationEntity})
	if err != nil {
		return nil, errors.WithMessagef(err, "application enrich with tokens")
	}

	return arr[0], nil
}

func (s Token) Revoke(ctx context.Context, req domain.TokenRevokeRequest) (*domain.ApplicationWithTokens, error) {
	app, err := s.appRep.GetApplicationById(ctx, req.AppId)
	if err != nil {
		return nil, errors.WithMessagef(err, "get application by id")
	}

	_, err = s.revokeTokens(ctx, req.Tokens)
	if err != nil {
		return nil, errors.WithMessagef(err, "reboke tokens")
	}

	res, err := s.appEnrich.EnrichWithTokens(ctx, []entity.Application{*app})
	if err != nil {
		return nil, errors.WithMessagef(err, "application enrich")
	}

	return res[0], nil
}

func (s Token) RevokeByAppId(ctx context.Context, appId int) (*domain.DeleteResponse, error) {
	tokens, err := s.tokenRep.GetTokenByAppIdList(ctx, []int{appId})
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
	err := s.tx.TokenRevokeTx(ctx, func(ctx context.Context, tx ITokenRevokeTx) error {
		deleted, err := tx.DeleteToken(ctx, tokens)
		if err != nil {
			return errors.WithMessagef(err, "tx delete token")
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
