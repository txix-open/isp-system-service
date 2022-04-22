package service

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"isp-system-service/domain"
	"isp-system-service/entity"
)

type IAccessListAccessListRep interface {
	GetAccessListByAppId(ctx context.Context, appId int) ([]entity.AccessList, error)
}

type IAccessListApplicationRep interface {
	GetApplicationById(ctx context.Context, id int) (*entity.Application, error)
}

type IAccessListSetOneTx interface {
	UpsertAccessList(ctx context.Context, e entity.AccessList) (int, error)
}

type IAccessListSetListTx interface {
	GetAccessListByAppId(ctx context.Context, appId int) ([]entity.AccessList, error)
	InsertArrayAccessList(ctx context.Context, entity []entity.AccessList) error
	DeleteAccessList(ctx context.Context, appId int, methods []string) error
	DeleteAccessListByAppId(ctx context.Context, appId int) ([]entity.AccessList, error)
}

type IAccessListTxRunner interface {
	AccessListSetOneTx(ctx context.Context, tx func(ctx context.Context, tx IAccessListSetOneTx) error) error
	AccessListSetListTx(ctx context.Context, tx func(ctx context.Context, tx IAccessListSetListTx) error) error
}

type IAccessListRedis interface {
	UpdateApplicationPermission(ctx context.Context, req entity.RedisApplicationPermission) error
	UpdateApplicationPermissionList(ctx context.Context, removed []string, added []interface{}) error
}

type AccessList struct {
	redis          IAccessListRedis
	tx             IAccessListTxRunner
	accessListRep  IAccessListAccessListRep
	applicationRep IAccessListApplicationRep
}

func NewAccessList(
	redis IAccessListRedis,
	tx IAccessListTxRunner,
	accessListRep IAccessListAccessListRep,
	applicationRep IAccessListApplicationRep,
) AccessList {
	return AccessList{
		redis:          redis,
		tx:             tx,
		accessListRep:  accessListRep,
		applicationRep: applicationRep,
	}
}

func (s AccessList) GetById(ctx context.Context, appId int) ([]domain.MethodInfo, error) {
	_, err := s.applicationRep.GetApplicationById(ctx, appId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	accessList, err := s.accessListRep.GetAccessListByAppId(ctx, appId)
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
	_, err := s.applicationRep.GetApplicationById(ctx, request.AppId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	var resp int
	err = s.tx.AccessListSetOneTx(ctx, func(ctx context.Context, tx IAccessListSetOneTx) error {
		resp, err = tx.UpsertAccessList(ctx, entity.AccessList{
			AppId:  request.AppId,
			Method: request.Method,
			Value:  request.Value,
		})
		if err != nil {
			return errors.WithMessagef(err, "upsert access list")
		}

		err = s.redis.UpdateApplicationPermission(ctx, entity.RedisApplicationPermission{
			AppId:  request.AppId,
			Method: request.Method,
			Value:  request.Value,
		})
		if err != nil {
			return errors.WithMessagef(err, "redis update application permission")
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
	_, err := s.applicationRep.GetApplicationById(ctx, req.AppId)
	if err != nil {
		return nil, errors.WithMessage(err, "get application by id")
	}

	err = s.tx.AccessListSetListTx(ctx, func(ctx context.Context, tx IAccessListSetListTx) error {
		oldAccessList := make([]entity.AccessList, 0)
		if req.RemoveOld {
			oldAccessList, err = tx.DeleteAccessListByAppId(ctx, req.AppId)
			if err != nil {
				return errors.WithMessage(err, "delete access list by app_id")
			}
		}

		newAccessList := make([]entity.AccessList, len(req.Methods))
		updateMethods := make([]string, len(req.Methods))
		redisMsetRequest := make([]interface{}, 0)
		for i, m := range req.Methods {
			updateMethods[i] = m.Method
			newAccessList[i] = entity.AccessList{
				AppId:  req.AppId,
				Method: m.Method,
				Value:  m.Value,
			}
			redisMsetRequest = append(redisMsetRequest, fmt.Sprintf("%d|%s", req.AppId, m.Method), m.Value)
		}

		redisDelRequest := make([]string, len(oldAccessList))
		for i, access := range oldAccessList {
			redisDelRequest[i] = fmt.Sprintf("%d|%s", req.AppId, access.Method)
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

		err = s.redis.UpdateApplicationPermissionList(ctx, redisDelRequest, redisMsetRequest)
		if err != nil {
			return errors.WithMessage(err, "redis update application permission list")
		}

		return nil
	})
	if err != nil {
		return nil, errors.WithMessage(err, "transaction access list set list")
	}

	accessList, err := s.accessListRep.GetAccessListByAppId(ctx, req.AppId)
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
