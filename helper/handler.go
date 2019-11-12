package helper

import (
	"isp-system-service/controller"
	"isp-system-service/domain"
	"isp-system-service/entity"

	"google.golang.org/grpc/metadata"
)

type SystemHandler struct {
	GetSystems         func([]int32) ([]entity.System, error)        `method:"get_systems" group:"system" inner:"true"`
	CreateUpdateSystem func(entity.System) (*entity.System, error)   `method:"create_update_system" group:"system" inner:"true"`
	GetSystemById      func(domain.Identity) (*entity.System, error) `method:"get_system_by_id" group:"system" inner:"true"`
	DeleteSystems      func([]int32) (domain.DeleteResponse, error)  `method:"delete_systems" group:"system" inner:"true"`
}

type DomainHandler struct {
	GetDomainsBySystemId func(metadata.MD) ([]entity.Domain, error)               `method:"get_domains_by_system_id" group:"domain" inner:"true"`
	CreateUpdateDomain   func(entity.Domain, metadata.MD) (*entity.Domain, error) `method:"create_update_domain" group:"domain" inner:"true"`
	GetDomainById        func(domain.Identity) (*entity.Domain, error)            `method:"get_domain_by_id" group:"domain" inner:"true"`
	DeleteDomains        func([]int32) (domain.DeleteResponse, error)             `method:"delete_domains" group:"domain" inner:"true"`
}

type ServiceHandler struct {
	GetService            func([]int32) ([]entity.Service, error)         `method:"get_service" group:"service" inner:"true"`
	GetServicesByDomainId func(domain.Identity) ([]entity.Service, error) `method:"get_services_by_domain_id" group:"service" inner:"true"`
	CreateUpdateService   func(entity.Service) (*entity.Service, error)   `method:"create_update_service" group:"service" inner:"true"`
	GetServiceById        func(domain.Identity) (*entity.Service, error)  `method:"get_service_by_id" group:"service" inner:"true"`
	DeleteService         func([]int32) (domain.DeleteResponse, error)    `method:"delete_service" group:"service" inner:"true"`
}

type ApplicationHandler struct {
	GetApplications            func([]int32) ([]*domain.AppWithToken, error)           `method:"get_applications" group:"application" inner:"true"`
	GetApplicationsByServiceId func(domain.Identity) ([]*domain.AppWithToken, error)   `method:"get_applications_by_service_id" group:"application" inner:"true"`
	CreateUpdateApplication    func(entity.Application) (*domain.AppWithToken, error)  `method:"create_update_application" group:"application" inner:"true"`
	GetApplicationById         func(domain.Identity) (*domain.AppWithToken, error)     `method:"get_application_by_id" group:"application" inner:"true"`
	DeleteApplications         func([]int32) (domain.DeleteResponse, error)            `method:"delete_applications" group:"application" inner:"true"`
	GetSystemTree              func(metadata.MD) ([]*domain.DomainWithServices, error) `method:"get_system_tree" group:"application" inner:"true"`
}

type TokenHandler struct {
	CreateToken        func(domain.CreateTokenRequest) (*domain.AppWithToken, error)  `method:"create_token" group:"token" inner:"true"`
	RevokeTokens       func(domain.RevokeTokensRequest) (*domain.AppWithToken, error) `method:"revoke_tokens" group:"token" inner:"true"`
	RevokeTokensForApp func(domain.Identity) (*domain.DeleteResponse, error)          `method:"revoke_tokens_for_app" group:"token" inner:"true"`
	GetTokensByAppId   func(domain.Identity) ([]entity.Token, error)                  `method:"get_tokens_by_app_id" group:"token" inner:"true"`
}

type AccessListHandler struct {
	GetById func(domain.Identity) (domain.ModuleMethods, error)        `method:"get_by_id" group:"access_list" inner:"true"`
	SetOne  func(entity.AccessList) (*domain.CountResponse, error)     `method:"set_one" group:"access_list" inner:"true"`
	SetList func(domain.SetListRequest) (*domain.CountResponse, error) `method:"set_list" group:"access_list" inner:"true"`
}

func GetSystemHandler() *SystemHandler {
	return &SystemHandler{
		GetSystems:         controller.System.GetSystems,
		CreateUpdateSystem: controller.System.CreateUpdateSystem,
		GetSystemById:      controller.System.GetSystemById,
		DeleteSystems:      controller.System.DeleteSystems,
	}
}

func GetDomainHandler() *DomainHandler {
	return &DomainHandler{
		GetDomainsBySystemId: controller.Domain.GetDomainsBySystemId,
		CreateUpdateDomain:   controller.Domain.CreateUpdateDomain,
		GetDomainById:        controller.Domain.GetDomainById,
		DeleteDomains:        controller.Domain.DeleteDomains,
	}
}

func GetServiceHandler() *ServiceHandler {
	return &ServiceHandler{
		GetService:            controller.Service.GetService,
		GetServicesByDomainId: controller.Service.GetServicesByDomainId,
		CreateUpdateService:   controller.Service.CreateUpdateService,
		GetServiceById:        controller.Service.GetServiceById,
		DeleteService:         controller.Service.DeleteService,
	}
}

func GetApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{
		GetApplications:            controller.Application.GetApplications,
		GetApplicationsByServiceId: controller.Application.GetApplicationsByServiceId,
		CreateUpdateApplication:    controller.Application.CreateUpdateApplication,
		GetApplicationById:         controller.Application.GetApplicationById,
		DeleteApplications:         controller.Application.DeleteApplications,
		GetSystemTree:              controller.Application.GetSystemTree,
	}
}

func GetTokenHandler() *TokenHandler {
	return &TokenHandler{
		CreateToken:        controller.Token.CreateToken,
		GetTokensByAppId:   controller.Token.GetTokensByAppId,
		RevokeTokens:       controller.Token.RevokeTokens,
		RevokeTokensForApp: controller.Token.RevokeTokensForApp,
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
