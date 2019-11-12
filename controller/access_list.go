package controller

import (
	"fmt"
	rd "github.com/go-redis/redis"
	rdLib "github.com/integration-system/isp-lib/redis"
	"isp-system-service/domain"
	"isp-system-service/entity"
	"isp-system-service/model"
	"isp-system-service/redis"
)

var AccessList accessListController

type accessListController struct{}

// GetById godoc
// @Tags accessList
// @Summary Получить список доступности методов для приложения
// @Description Возвращает список методов для приложения, для которых заданы настройки доступа
// @Accept  json
// @Produce  json
// @Param body body domain.Identity false "идентификатор приложения"
// @Success 200 {array} domain.ModuleMethods "список методов по названию модуля"
// @Failure 500 {object} structure.GrpcError
// @Router /access_list/get_by_id [POST]
func (a accessListController) GetById(req domain.Identity) (domain.ModuleMethods, error) {
	if accessList, err := model.AccessListRep.GetByAppId(req.Id); err != nil {
		return nil, err
	} else {
		return a.convertAccessList(accessList), nil
	}
}

// SetOne godoc
// @Tags accessList
// @Summary Настроить доступность метода для приложения
// @Description Настраивает достуность
// @Accept  json
// @Produce  json
// @Param body body entity.AccessList false "объект для настройки доступа"
// @Success 200 {object} domain.CountResponse "количество измененных строк"
// @Failure 500 {object} structure.GrpcError
// @Router /access_list/set_one [POST]
func (c accessListController) SetOne(request entity.AccessList) (*domain.CountResponse, error) {
	var (
		resp = 0
		err  error
	)
	if err := model.DbClient.RunInTransaction(func(repository model.AccessListRepository) error {
		if resp, err = repository.Upsert(request); err != nil {
			return err
		}

		if _, err := redis.Client.Get().UseDb(rdLib.ApplicationPermissionDb, func(p rd.Pipeliner) error {
			key := fmt.Sprintf("%d|%s", request.AppId, request.Method)
			if _, err := p.Set(key, request.Value, 0).Result(); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		} else {
			return nil
		}
	}); err != nil {
		return nil, err
	} else {
		return &domain.CountResponse{Count: resp}, nil
	}
}

// SetList godoc
// @Tags accessList
// @Summary Настройть доступность списка методов для приложения
// @Description Возвращает массив приложений с токенами по их идентификаторам
// @Accept  json
// @Produce  json
// @Param body body domain.SetListRequest false "объект настройки доступа"
// @Success 200 {object} domain.CountResponse "количество добавленных строк"
// @Failure 500 {object} structure.GrpcError
// @Router /access_list/set_list [POST]
func (c accessListController) SetList(request domain.SetListRequest) (*domain.CountResponse, error) {
	resp := 0
	if err := model.DbClient.RunInTransaction(func(repository model.AccessListRepository) error {

		oldAccessList, err := repository.DeleteById(request.AppId)
		if err != nil {
			return err
		}

		insertRequest := make([]entity.AccessList, len(request.Methods))
		redisMsetRequest := make([]interface{}, 0)
		for i, m := range request.Methods {
			insertRequest[i] = entity.AccessList{
				AppId:  request.AppId,
				Method: m.Method,
				Value:  m.Value,
			}
			redisMsetRequest = append(redisMsetRequest, fmt.Sprintf("%d|%s", request.AppId, m.Method), m.Value)
		}

		redisDelRequest := make([]string, len(oldAccessList))
		for i, access := range oldAccessList {
			redisDelRequest[i] = fmt.Sprintf("%d|%s", request.AppId, access.Method)
		}

		if resp, err = repository.Insert(insertRequest); err != nil {
			return err
		}

		if _, err := redis.Client.Get().UseDb(rdLib.ApplicationPermissionDb, func(p rd.Pipeliner) error {
			if len(redisDelRequest) > 0 {
				if _, err := p.Del(redisDelRequest...).Result(); err != nil {
					return err
				}
			}

			if _, err := p.MSet(redisMsetRequest...).Result(); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		} else {
			return nil
		}
	}); err != nil {
		return nil, err
	} else {
		return &domain.CountResponse{Count: resp}, nil
	}
}

func (c accessListController) convertAccessList(accessLists []entity.AccessList) domain.ModuleMethods {
	response := make(map[string][]domain.MethodInfo)
	for _, access := range accessLists {
		moduleName, method := c.extractModuleName(access.Method)
		if methodStore, ok := response[moduleName]; ok {
			response[moduleName] = append(methodStore, domain.MethodInfo{Value: access.Value, Method: method})
		} else {
			response[moduleName] = []domain.MethodInfo{{Value: access.Value, Method: method}}
		}
	}
	return response
}

func (accessListController) extractModuleName(method string) (string, string) {
	firstFound := false
	for i, value := range method {
		if value == '/' {
			if firstFound {
				return method[:i], method[i:]
			} else {
				firstFound = true
			}
		}
	}
	return method, ""
}
