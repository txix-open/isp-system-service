package controller

import (
	"context"

	"github.com/pkg/errors"
	"isp-system-service/domain"
)

type SecureService interface {
	Authenticate(ctx context.Context, token string) (*domain.AuthData, error)
	Authorize(ctx context.Context, appId int, endpoint string) (bool, error)
}

type Secure struct {
	service SecureService
}

func NewSecure(service SecureService) Secure {
	return Secure{
		service: service,
	}
}

// Authenticate godoc
//
//	@Tags			secure
//	@Summary		Метод аутентификации токена
//	@Description	Проверяет наличие токена в системе,
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.AuthenticateRequest	true	"Тело запроса"
//	@Success		200		{array}		domain.AuthenticateResponse
//	@Failure		500		{object}	domain.GrpcError
//	@Router			/secure/authenticate [POST]
func (c Secure) Authenticate(ctx context.Context, req domain.AuthenticateRequest) (*domain.AuthenticateResponse, error) {
	result, err := c.service.Authenticate(ctx, req.Token)
	switch {
	case errors.Is(err, domain.ErrTokenNotFound):
		return &domain.AuthenticateResponse{
			Authenticated: false,
			ErrorReason:   domain.ErrTokenNotFound.Error(),
		}, nil
	case errors.Is(err, domain.ErrTokenExpired):
		return &domain.AuthenticateResponse{
			Authenticated: false,
			ErrorReason:   domain.ErrTokenExpired.Error(),
		}, nil
	case err != nil:
		return nil, errors.WithMessage(err, "authenticate")
	default:
		return &domain.AuthenticateResponse{
			Authenticated: true,
			AuthData:      result,
		}, nil
	}
}

// Authorize godoc
//
//	@Tags			secure
//	@Summary		Метод авторизации приложения
//	@Description	Проверяет доступ приложения к запрашиваемому ендпоинту
//	@Accept			json
//	@Produce		json
//	@Param			body	body		domain.AuthorizeRequest	true	"Тело запрос"
//	@Success		200		{array}		domain.AuthorizeResponse
//	@Failure		500		{object}	domain.GrpcError
//	@Router			/secure/authorize [POST]
func (c Secure) Authorize(ctx context.Context, req domain.AuthorizeRequest) (*domain.AuthorizeResponse, error) {
	result, err := c.service.Authorize(ctx, req.ApplicationId, req.Endpoint)
	switch {
	case errors.Is(err, domain.ErrAccessListNotFound):
		return &domain.AuthorizeResponse{
			Authorized: false,
		}, nil
	case err != nil:
		return nil, errors.WithMessage(err, "authorize")
	default:
		return &domain.AuthorizeResponse{
			Authorized: result,
		}, nil
	}
}
