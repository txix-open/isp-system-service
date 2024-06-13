package secure

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"isp-system-service/domain"
	"isp-system-service/entity"
)

type TokenRep interface {
	AuthDataByToken(ctx context.Context, token string) (*entity.AuthData, error)
}

type AccessListRep interface {
	GetAccessListByAppIdAndMethod(ctx context.Context, appId int, method string) (*entity.AccessList, error)
}

type ApplicationGroup struct {
	tokenRep      TokenRep
	accessListRep AccessListRep
}

func NewApplicationGroup(
	tokenRep TokenRep,
	accessListRep AccessListRep,
) ApplicationGroup {
	return ApplicationGroup{
		tokenRep:      tokenRep,
		accessListRep: accessListRep,
	}
}

func (s ApplicationGroup) Authenticate(ctx context.Context, token string) (*domain.AuthData, error) {
	authData, err := s.tokenRep.AuthDataByToken(ctx, token)
	if err != nil {
		return nil, errors.WithMessage(err, "get auth data by token")
	}

	if authData.ExpireTime != -1 &&
		authData.CreatedAt.Add(time.Millisecond*time.Duration(authData.ExpireTime)).Before(time.Now().UTC()) {
		return nil, domain.ErrTokenExpired
	}

	return &domain.AuthData{
		SystemId:           authData.SystemId,
		DomainId:           authData.DomainId,
		ApplicationGroupId: authData.ApplicationGroupIdId,
		ApplicationId:      authData.AppId,
	}, nil
}

func (s ApplicationGroup) Authorize(ctx context.Context, appId int, endpoint string) (bool, error) {
	accessList, err := s.accessListRep.GetAccessListByAppIdAndMethod(ctx, appId, endpoint)
	if err != nil {
		return false, errors.WithMessage(err, "get access list by app_id and method")
	}

	return accessList.Value, nil
}
