package controller

import (
	"context"
	"fmt"

	"isp-system-service/domain"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/grpc/apierrors"
	"google.golang.org/grpc/codes"
)

type TokenService interface {
	GetByAppId(ctx context.Context, appId int) ([]domain.Token, error)
	Create(ctx context.Context, req domain.TokenCreateRequest) (*domain.ApplicationWithTokens, error)
	Revoke(ctx context.Context, req domain.TokenRevokeRequest) (*domain.ApplicationWithTokens, error)
	RevokeByAppId(ctx context.Context, appId int) (*domain.DeleteResponse, error)
}

type Token struct {
	service TokenService
}

func NewToken(service TokenService) Token {
	return Token{
		service: service,
	}
}

// GetByAppId godoc
//
//	@Tags			token
//	@Summary		Получить токены по идентификатору приложения
//	@Description	Возвращает список токенов, привязанных к приложению
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.Identity	true	"Идентификатор приложения"
//	@Success		200		{array}		entity.Token
//	@Failure		500		{object}	apierrors.Error
//	@Router			/token/get_tokens_by_app_id [POST]
func (c Token) GetByAppId(ctx context.Context, req domain.Identity) ([]domain.Token, error) {
	return c.service.GetByAppId(ctx, req.Id)
}

// Create godoc
//
//	@Tags			token
//	@Summary		Создать токен
//	@Description	Создает токен и привязывает его к приложению
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.TokenCreateRequest	true	"Объект создания токена"
//	@Success		200		{object}	domain.ApplicationWithTokens
//	@Failure		500		{object}	apierrors.Error
//	@Router			/token/create_token [POST]
func (c Token) Create(ctx context.Context, req domain.TokenCreateRequest) (*domain.ApplicationWithTokens, error) {
	result, err := c.service.Create(ctx, req)
	switch {
	case errors.Is(err, domain.ErrApplicationNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeApplicationNotFound,
			fmt.Sprintf("application with id %d not found", req.AppId),
			err,
		)
	case errors.Is(err, domain.ErrAppGroupNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeAppGroupNotFound,
			fmt.Sprintf("service for app_id id %d not found", req.AppId),
			err,
		)
	case errors.Is(err, domain.ErrDomainNotFound):
		return nil, apierrors.New(
			codes.NotFound,
			domain.ErrCodeDomainNotFound,
			fmt.Sprintf("domain for app_id id %d not found", req.AppId),
			err,
		)
	case err != nil:
		return nil, err
	default:
		return result, nil
	}
}

// Revoke godoc
//
//	@Tags			token
//	@Summary		Отозвать токены
//	@Description	Отвязывает токены от приложений и удаляет их
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.TokenRevokeRequest	true	"Объект отзыва токенов"
//	@Success		200		{object}	domain.ApplicationWithTokens
//	@Failure		500		{object}	apierrors.Error
//	@Router			/token/revoke_tokens [POST]
func (c Token) Revoke(ctx context.Context, req domain.TokenRevokeRequest) (*domain.ApplicationWithTokens, error) {
	return c.service.Revoke(ctx, req)
}

// RevokeForApp godoc
//
//	@Tags			token
//	@Summary		Отозвать токены для приложения
//	@Description	Отвязывает токены от приложений и удаляет их по идентификатору приложения
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.Identity	true	"Идентификатор приложения"
//	@Success		200		{object}	domain.DeleteResponse
//	@Failure		500		{object}	apierrors.Error
//	@Router			/token/revoke_tokens_for_app [POST]
func (c Token) RevokeForApp(ctx context.Context, req domain.Identity) (*domain.DeleteResponse, error) {
	return c.service.RevokeByAppId(ctx, req.Id)
}
