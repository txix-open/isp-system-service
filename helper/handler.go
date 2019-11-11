package helper

import (
	"isp-system-service/controller"
	"isp-system-service/entity"

	"google.golang.org/grpc/metadata"
)

type SystemHandler struct {
	GetSystems         func([]int32) ([]entity.System, error)            `method:"get_systems" group:"system" inner:"true"`
	CreateUpdateSystem func(entity.System) (*entity.System, error)       `method:"create_update_system" group:"system" inner:"true"`
	GetSystemById      func(controller.Identity) (*entity.System, error) `method:"get_system_by_id" group:"system" inner:"true"`
	DeleteSystems      func([]int32) (controller.DeleteResponse, error)  `method:"delete_systems" group:"system" inner:"true"`
}

type DomainHandler struct {
	//GetDomains           func([]int32) ([]entity.Domain, error)                 `method:"get_domains" group:"domain" inner:"true"`
	GetDomainsBySystemId func(metadata.MD) ([]entity.Domain, error)               `method:"get_domains_by_system_id" group:"domain" inner:"true"`
	CreateUpdateDomain   func(entity.Domain, metadata.MD) (*entity.Domain, error) `method:"create_update_domain" group:"domain" inner:"true"`
	GetDomainById        func(controller.Identity) (*entity.Domain, error)        `method:"get_domain_by_id" group:"domain" inner:"true"`
	DeleteDomains        func([]int32) (controller.DeleteResponse, error)         `method:"delete_domains" group:"domain" inner:"true"`
}

type ServiceHandler struct {
	GetService            func([]int32) ([]entity.Service, error)             `method:"get_service" group:"service" inner:"true"`
	GetServicesByDomainId func(controller.Identity) ([]entity.Service, error) `method:"get_services_by_domain_id" group:"service" inner:"true"`
	CreateUpdateService   func(entity.Service) (*entity.Service, error)       `method:"create_update_service" group:"service" inner:"true"`
	GetServiceById        func(controller.Identity) (*entity.Service, error)  `method:"get_service_by_id" group:"service" inner:"true"`
	DeleteService         func([]int32) (controller.DeleteResponse, error)    `method:"delete_service" group:"service" inner:"true"`
}

type ApplicationHandler struct {
	GetApplications            func([]int32) ([]*controller.AppWithToken, error)             `method:"get_applications" group:"application" inner:"true"`
	GetApplicationsByServiceId func(controller.Identity) ([]*controller.AppWithToken, error) `method:"get_applications_by_service_id" group:"application" inner:"true"`
	CreateUpdateApplication    func(entity.Application) (*controller.AppWithToken, error)    `method:"create_update_application" group:"application" inner:"true"`
	GetApplicationById         func(controller.Identity) (*controller.AppWithToken, error)   `method:"get_application_by_id" group:"application" inner:"true"`
	DeleteApplications         func([]int32) (controller.DeleteResponse, error)              `method:"delete_applications" group:"application" inner:"true"`
	GetSystemTree              func(metadata.MD) ([]*controller.DomainWithServices, error)   `method:"get_system_tree" group:"application" inner:"true"`
}

type TokenHandler struct {
	CreateToken        func(controller.CreateTokenRequest) (*controller.AppWithToken, error)  `method:"create_token" group:"token" inner:"true"`
	RevokeTokens       func(controller.RevokeTokensRequest) (*controller.AppWithToken, error) `method:"revoke_tokens" group:"token" inner:"true"`
	RevokeTokensForApp func(controller.Identity) (*controller.DeleteResponse, error)          `method:"revoke_tokens_for_app" group:"token" inner:"true"`
	GetTokensByAppId   func(controller.Identity) ([]entity.Token, error)                      `method:"get_tokens_by_app_id" group:"token" inner:"true"`
}

type AccessListHandler struct {
	GetById func(controller.Identity) (controller.ModuleMethods, error)        `method:"get_by_id" group:"access_list" inner:"true"`
	SetOne  func(entity.AccessList) (*controller.CountResponse, error)         `method:"set_one" group:"access_list" inner:"true"`
	SetList func(controller.SetListRequest) (*controller.CountResponse, error) `method:"set_list" group:"access_list" inner:"true"`
}

func GetSystemHandler() *SystemHandler {
	return &SystemHandler{
		GetSystems:         controller.GetSystems,
		CreateUpdateSystem: controller.CreateUpdateSystem,
		GetSystemById:      controller.GetSystemById,
		DeleteSystems:      controller.DeleteSystems,
	}
}

func GetDomainHandler() *DomainHandler {
	return &DomainHandler{
		//GetDomains:           controller.GetDomains,
		GetDomainsBySystemId: controller.GetDomainsBySystemId,
		CreateUpdateDomain:   controller.CreateUpdateDomain,
		GetDomainById:        controller.GetDomainById,
		DeleteDomains:        controller.DeleteDomains,
	}
}

func GetServiceHandler() *ServiceHandler {
	return &ServiceHandler{
		GetService:            controller.GetService,
		GetServicesByDomainId: controller.GetServicesByDomainId,
		CreateUpdateService:   controller.CreateUpdateService,
		GetServiceById:        controller.GetServiceById,
		DeleteService:         controller.DeleteService,
	}
}

func GetApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{
		GetApplications:            controller.GetApplications,
		GetApplicationsByServiceId: controller.GetApplicationsByServiceId,
		CreateUpdateApplication:    controller.CreateUpdateApplication,
		GetApplicationById:         controller.GetApplicationById,
		DeleteApplications:         controller.DeleteApplications,
		GetSystemTree:              controller.GetSystemTree,
	}
}

func GetTokenHandler() *TokenHandler {
	return &TokenHandler{
		CreateToken:        controller.CreateToken,
		GetTokensByAppId:   controller.GetTokensByAppId,
		RevokeTokens:       controller.RevokeTokens,
		RevokeTokensForApp: controller.RevokeTokensForApp,
	}
}

func GetAccessListHandler() *AccessListHandler {
	return &AccessListHandler{
		GetById: controller.AccessList.GetById,
		SetOne:  controller.AccessList.SetOne,
		SetList: controller.AccessList.SetList,
	}
}

func GetAllHandlers() []interface{} {
	return []interface{}{
		GetSystemHandler(),
		GetDomainHandler(),
		GetServiceHandler(),
		GetApplicationHandler(),
		GetTokenHandler(),
		GetAccessListHandler(),
	}
}
