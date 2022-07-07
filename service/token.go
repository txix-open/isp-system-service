package service

import (
	"context"

	"github.com/pkg/errors"
	"isp-system-service/domain"
	"isp-system-service/entity"
)

type ITokenJwt interface {
	CreateApplicationToken(appId int, expireTime int) (string, error)
}

type ITokenAppEnrich interface {
	EnrichWithTokens(ctx context.Context, apps []entity.Application) ([]*domain.ApplicationWithTokens, error)
}

type ITokenApplicationRep interface {
	GetApplicationById(ctx context.Context, id int) (*entity.Application, error)
}

type ITokenServiceRep interface {
	GetServiceById(ctx context.Context, id int) (*entity.Service, error)
}

type ITokenDomainRep interface {
	GetDomainById(ctx context.Context, id int) (*entity.Domain, error)
}

type ITokenTokenRep interface {
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

type ITokenRedis interface {
	SetApplicationToken(ctx context.Context, req entity.RedisSetToken) error
	DeleteToken(ctx context.Context, tokens []string) error
}

type Token struct {
	redisRep          ITokenRedis
	defaultExpireTime int
	jwt               ITokenJwt
	appEnrich         ITokenAppEnrich
	tx                ITokenTxRunner
	appRep            ITokenApplicationRep
	domainRep         ITokenDomainRep
	serviceRep        ITokenServiceRep
	tokenRep          ITokenTokenRep
}

func NewToken(
	redisRep ITokenRedis,
	defaultExpireTime int,
	jwtGenerate ITokenJwt,
	appEnrich ITokenAppEnrich,
	tx ITokenTxRunner,
	appRep ITokenApplicationRep,
	domainRep ITokenDomainRep,
	serviceRep ITokenServiceRep,
	tokenRep ITokenTokenRep,
) Token {
	return Token{
		redisRep:          redisRep,
		defaultExpireTime: defaultExpireTime,
		appEnrich:         appEnrich,
		jwt:               jwtGenerate,
		tx:                tx,
		appRep:            appRep,
		domainRep:         domainRep,
		serviceRep:        serviceRep,
		tokenRep:          tokenRep,
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

	domainEntity, err := s.domainRep.GetDomainById(ctx, serviceEntity.DomainId)
	if err != nil {
		return nil, errors.WithMessagef(err, "get domain by id")
	}

	expTime := req.ExpireTimeMs
	if expTime == 0 {
		expTime = s.defaultExpireTime
	}

	token, err := s.jwt.CreateApplicationToken(req.AppId, expTime)
	if err != nil {
		return nil, errors.WithMessagef(err, "jwt create application token")
	}

	err = s.tx.TokenCreateTx(ctx, func(ctx context.Context, tx ITokenCreateTx) error {
		tokenEntity, err := tx.SaveToken(ctx, token, req.AppId, expTime)
		if err != nil {
			return errors.WithMessagef(err, "tx save token")
		}

		err = s.redisRep.SetApplicationToken(ctx, entity.RedisSetToken{
			Token:               tokenEntity.Token,
			ExpireTime:          tokenEntity.ExpireTime,
			DomainIdentity:      domainEntity.Id,
			ServiceIdentity:     serviceEntity.Id,
			ApplicationIdentity: applicationEntity.Id,
		})
		if err != nil {
			return errors.WithMessage(err, "redis set token")
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

		err = s.redisRep.DeleteToken(ctx, tokens)
		if err != nil {
			return errors.WithMessage(err, "redis delete token")
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