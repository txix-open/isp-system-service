package controller

import (
	"isp-system-service/entity"
	"isp-system-service/model"

	"github.com/integration-system/isp-lib/database"
	"github.com/integration-system/isp-lib/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GetApplications(list []int32) ([]*AppWithToken, error) {
	res, err := model.AppRep.GetApplications(list)
	if err != nil {
		return nil, err
	}
	return enrichWithTokens(res...)
}

func GetApplicationsByServiceId(identity Identity) ([]*AppWithToken, error) {
	arr, err := model.AppRep.GetApplicationsByServiceId(identity.Id)
	if err != nil {
		return nil, err
	}
	return enrichWithTokens(arr...)
}

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
		dws := &DomainWithServices{Id: d.Id, Name: d.Name, Description: d.Description}
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
		swa := &ServiceWithApps{Id: s.Id, Name: s.Name, Description: s.Description}
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
		awt := &AppWithToken{App: a}
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
