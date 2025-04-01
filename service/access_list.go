package service

import (
	"context"

	"isp-system-service/domain"
	"isp-system-service/entity"

	"github.com/pkg/errors"
)

type AccessListRepo interface {
	GetAccessListByAppId(ctx context.Context, appId int) ([]entity.AccessList, error)
}

type AccessListSetOneTx interface {
	UpsertAccessList(ctx context.Context, e entity.AccessList) (int, error)
}

type AccessListSetListTx interface {
	GetAccessListByAppId(ctx context.Context, appId int) ([]entity.AccessList, error)
	InsertArrayAccessList(ctx context.Context, entity []entity.AccessList) error
	DeleteAccessList(ctx context.Context, appId int, methods []string) error
	DeleteAccessListByAppId(ctx context.Context, appId int) ([]entity.AccessList, error)
}

type AccessListTxRunner interface {
	AccessListSetOneTx(ctx context.Context, tx func(ctx context.Context, tx AccessListSetOneTx) error) error
	AccessListSetListTx(ctx context.Context, tx func(ctx context.Context, tx AccessListSetListTx) error) error
}

type AccessList struct {
	tx             AccessListTxRunner
	accessListRepo AccessListRepo
	appRepo        ApplicationRepo
}

func NewAccessList(
	tx AccessListTxRunner,
	accessListRepo AccessListRepo,
	appRepo ApplicationRepo,
) AccessList {
	return AccessList{
		tx:             tx,
		accessListRepo: accessListRepo,
		appRepo:        appRepo,
	}
}

func (s AccessList) GetById(ctx context.Context, appId int) ([]domain.MethodInfo, error) {
	_, err := s.appRepo.GetApplicationById(ctx, appId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	accessList, err := s.accessListRepo.GetAccessListByAppId(ctx, appId)
	if err != nil {
		return nil, errors.WithMessagef(err, "get access list by app_id")
	}

	methodInfos := make([]domain.MethodInfo, len(accessList))
	for i, access := range accessList {
		methodInfos[i] = domain.MethodInfo{
			Method: access.Method,
			Value:  access.Value,
		}
	}

	return methodInfos, nil
}

func (s AccessList) SetOne(ctx context.Context, request domain.AccessListSetOneRequest) (*domain.AccessListSetOneResponse, error) {
	_, err := s.appRepo.GetApplicationById(ctx, request.AppId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	var resp int
	err = s.tx.AccessListSetOneTx(ctx, func(ctx context.Context, tx AccessListSetOneTx) error {
		resp, err = tx.UpsertAccessList(ctx, entity.AccessList{
			AppId:  request.AppId,
			Method: request.Method,
			Value:  request.Value,
		})
		if err != nil {
			return errors.WithMessagef(err, "upsert access list")
		}

		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "transaction access list set one")
	}

	return &domain.AccessListSetOneResponse{
		Count: resp,
	}, nil
}

func (s AccessList) SetList(ctx context.Context, req domain.AccessListSetListRequest) ([]domain.MethodInfo, error) {
	_, err := s.appRepo.GetApplicationById(ctx, req.AppId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	err = s.tx.AccessListSetListTx(ctx, func(ctx context.Context, tx AccessListSetListTx) error {
		if req.RemoveOld {
			_, err = tx.DeleteAccessListByAppId(ctx, req.AppId)
			if err != nil {
				return errors.WithMessage(err, "delete access list by app_id")
			}
		}

		newAccessList := make([]entity.AccessList, len(req.Methods))
		updateMethods := make([]string, len(req.Methods))
		for i, m := range req.Methods {
			updateMethods[i] = m.Method
			newAccessList[i] = entity.AccessList{
				AppId:  req.AppId,
				Method: m.Method,
				Value:  m.Value,
			}
		}

		if len(newAccessList) > 0 {
			err = tx.DeleteAccessList(ctx, req.AppId, updateMethods)
			if err != nil {
				return errors.WithMessage(err, "delete access_list")
			}

			err = tx.InsertArrayAccessList(ctx, newAccessList)
			if err != nil {
				return errors.WithMessage(err, "insert access_list")
			}
		}

		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "transaction access list set list")
	}

	accessList, err := s.accessListRepo.GetAccessListByAppId(ctx, req.AppId)
	if err != nil {
		return nil, errors.WithMessagef(err, "get access list by app_id")
	}

	methodInfos := make([]domain.MethodInfo, len(accessList))
	for i, access := range accessList {
		methodInfos[i] = domain.MethodInfo{
			Method: access.Method,
			Value:  access.Value,
		}
	}

	return methodInfos, nil
}
