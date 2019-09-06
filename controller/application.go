package controller

import (
	"isp-system-service/entity"
	"isp-system-service/model"

	"github.com/integration-system/isp-lib/database"
	_ "github.com/integration-system/isp-lib/structure"
	"github.com/integration-system/isp-lib/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetApplications godoc
// @Tags application
// @Summary Получить список приложений
// @Description Возвращает массив приложений с токенами по их идентификаторам
// @Accept  json
// @Produce  json
// @Param body body []integer false "Массив идентификаторов приложений"
// @Success 200 {array} controller.AppWithToken
// @Failure 500 {object} structure.GrpcError
// @Router /application/get_applications [POST]
func GetApplications(list []int32) ([]*AppWithToken, error) {
	res, err := model.AppRep.GetApplications(list)
	if err != nil {
		return nil, err
	}
	return enrichWithTokens(res...)
}

// GetApplicationsByServiceId godoc
// @Tags application
// @Summary Получить список приложений по идентификатору сервиса
// @Description Возвращает список приложений по запрошенныму идентификатору сервиса
// @Accept  json
// @Produce  json
// @Param body body controller.Identity true "Идентификатор серсиса"
// @Success 200 {array} controller.AppWithToken
// @Failure 500 {object} structure.GrpcError
// @Router /application/get_applications_by_service_id [POST]
func GetApplicationsByServiceId(identity Identity) ([]*AppWithToken, error) {
	arr, err := model.AppRep.GetApplicationsByServiceId(identity.Id)
	if err != nil {
		return nil, err
	}
	return enrichWithTokens(arr...)
}

// CreateUpdateApplication godoc
// @Tags application
// @Summary Создать/обновить приложение
// @Description Если приложение с такими идентификатором существует, то обновляет данные, если нет, то добавляет данные в базу
// @Accept  json
// @Produce  json
// @Param body body entity.Application true "Объект приложения"
// @Success 200 {object} controller.AppWithToken
// @Failure 400 {object} structure.GrpcError
// @Failure 404 {object} structure.GrpcError
// @Failure 409 {object} structure.GrpcError
// @Failure 500 {object} structure.GrpcError
// @Router /application/create_update_application [POST]
func CreateUpdateApplication(app entity.Application) (*AppWithToken, error) {
	existed, err := model.AppRep.GetApplicationByNameAndServiceId(app.Name, app.ServiceId)
	if err != nil {
		return nil, err
	}
	s, e := model.ServiceRep.GetServiceById(app.ServiceId)
	if e != nil {
		return nil, e
	}
	if s == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Service with id %d not found", app.ServiceId)
	}
	if app.Id == 0 {
		if existed != nil {
			return nil, status.Errorf(codes.AlreadyExists, "Application with name %s already exists", app.Name)
		}
		app, e := model.AppRep.CreateApplication(app)
		enrichWithTokens(app)
		return &AppWithToken{App: app, Tokens: []entity.Token{}}, e
	} else {
		if existed != nil && existed.Id != app.Id {
			return nil, status.Errorf(codes.AlreadyExists, "Application with name %s already exists", app.Name)
		}
		existed, err = model.AppRep.GetApplicationById(app.Id)
		if err != nil {
			return nil, err
		}
		if existed == nil {
			return nil, status.Errorf(codes.NotFound, "Application with id %d not found", app.Id)
		}
		app, e := model.AppRep.UpdateApplication(app)
		if e != nil {
			return nil, e
		}
		arr, err := enrichWithTokens(app)
		if err != nil {
			return nil, err
		}
		return arr[0], nil
	}
}

// GetApplicationById godoc
// @Tags application
// @Summary Получить приложение по идентификатору
// @Description  Возвращает описание приложения по его идентификатору
// @Accept  json
// @Produce  json
// @Param body body controller.Identity true "Идентификатор приложения"
// @Success 200 {object} controller.AppWithToken
// @Failure 400 {object} structure.GrpcError
// @Failure 404 {object} structure.GrpcError
// @Failure 500 {object} structure.GrpcError
// @Router /application/get_application_by_id [POST]
func GetApplicationById(identity Identity) (*AppWithToken, error) {
	application, err := model.AppRep.GetApplicationById(identity.Id)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, status.Errorf(codes.NotFound, "Application with id %d not found", identity.Id)
	}
	arr, err := enrichWithTokens(*application)
	if err != nil {
		return nil, err
	}
	return arr[0], nil
}

// DeleteApplications godoc
// @Tags application
// @Summary Удалить приложения
// @Description Удаляет приложения по списку их идентификаторов, возвращает количество удаленных приложений
// @Accept  json
// @Produce  json
// @Param body body []integer false "Массив идентификаторов приложений"
// @Success 200 {object} controller.AppWithToken
// @Failure 400 {object} structure.GrpcError
// @Failure 500 {object} structure.GrpcError
// @Router /application/delete_applications [POST]
func DeleteApplications(list []int32) (DeleteResponse, error) {
	if len(list) == 0 {
		return DeleteResponse{}, status.Errorf(codes.InvalidArgument, "At least one id are required")
	}
	var count = 0
	err := database.RunInTransaction(func(appRep model.AppRepository, tokenRep model.TokenRepository) error {
		for _, appId := range list {
			_, err := revokeTokensForApp(Identity{appId}, &tokenRep)
			if err != nil {
				return err
			}
		}
		res, err := appRep.DeleteApplications(list)
		if err != nil {
			return err
		}
		count = res
		return nil
	})
	if err != nil {
		return DeleteResponse{}, err
	}
	return DeleteResponse{Deleted: count}, nil
}

// GetSystemTree godoc
// @Tags application
// @Summary Метод получения системного дерева
// @Description Возвращает описание взаимосвязей сервисов и приложений
// @Accept  json
// @Produce  json
// @Success 200 {array} controller.DomainWithServices
// @Failure 500 {object} structure.GrpcError
// @Router /application/get_system_tree [POST]
func GetSystemTree(md metadata.MD) ([]*DomainWithServices, error) {
	sysId, err := utils.ResolveMetadataIdentity(utils.SystemIdHeader, md)
	if err != nil {
		return nil, err
	}
	domains, err := model.DomainRep.GetDomainsBySystemId(int32(sysId))
	if err != nil {
		return nil, err
	}

	l := len(domains)
	if l == 0 {
		return []*DomainWithServices{}, nil
	}
	idList := make([]int32, l)
	domainsMap := make(map[int32]*DomainWithServices, l)
	res := make([]*DomainWithServices, l)
	for i, d := range domains {
		idList[i] = d.Id
		dws := &DomainWithServices{Id: d.Id, Name: d.Name, Description: d.Description, Services: make([]*ServiceWithApps, 0)}
		domainsMap[d.Id] = dws
		res[i] = dws
	}
	services, err := model.ServiceRep.GetServicesByDomainId(idList...)
	if err != nil {
		return nil, err
	}

	l = len(services)
	idList = make([]int32, l)
	servicesMap := make(map[int32]*ServiceWithApps, l)
	for i, s := range services {
		idList[i] = s.Id
		d := domainsMap[s.DomainId]
		swa := &ServiceWithApps{Id: s.Id, Name: s.Name, Description: s.Description, Apps: make([]*SimpleApp, 0)}
		d.Services = append(d.Services, swa)
		servicesMap[s.Id] = swa
	}
	apps, err := model.AppRep.GetApplicationsByServiceId(idList...)
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		s := servicesMap[app.ServiceId]
		s.Apps = append(s.Apps, &SimpleApp{Id: app.Id,
			Name:        app.Name,
			Type:        app.Type,
			Description: app.Description,
			Tokens:      make([]entity.Token, 0),
		})
	}

	return res, nil
}

func enrichWithTokens(apps ...entity.Application) ([]*AppWithToken, error) {
	l := len(apps)
	if len(apps) == 0 {
		return []*AppWithToken{}, nil
	}
	idList := make([]int32, l)
	appMap := make(map[int32]*AppWithToken, l)
	enriched := make([]*AppWithToken, l)
	for i, a := range apps {
		idList[i] = a.Id
		awt := &AppWithToken{App: a, Tokens: make([]entity.Token, 0)}
		appMap[a.Id] = awt
		enriched[i] = awt
	}
	tokens, err := model.TokenRep.GetTokensByAppId(idList...)
	if err != nil {
		return nil, err
	}
	for _, v := range tokens {
		app := appMap[v.AppId]
		app.Tokens = append(app.Tokens, v)
	}
	return enriched, nil

}
