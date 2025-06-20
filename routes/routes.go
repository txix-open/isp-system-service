package routes

import (
	"isp-system-service/controller"

	"github.com/txix-open/isp-kit/cluster"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/endpoint"
)

type Controllers struct {
	AccessList  controller.AccessList
	Domain      controller.Domain
	Service     controller.Service
	Application controller.Application
	AppGroup    controller.AppGroup
	Token       controller.Token
	Secure      controller.Secure
}

func EndpointDescriptors() []cluster.EndpointDescriptor {
	return endpointDescriptors(Controllers{})
}

func Handler(wrapper endpoint.Wrapper, c Controllers) *grpc.Mux {
	muxer := grpc.NewMux()
	for _, descriptor := range endpointDescriptors(c) {
		muxer.Handle(descriptor.Path, wrapper.Endpoint(descriptor.Handler))
	}
	return muxer
}

func endpointDescriptors(c Controllers) []cluster.EndpointDescriptor {
	return concatCluster(
		secureCluster(c),
		accessListCluster(c),
		domainCluster(c),
		serviceCluster(c),
		applicationCluster(c),
		tokenCluster(c),
		applicationGroupCluster(c),
	)
}

func secureCluster(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{
		{
			Path:    "system/secure/authenticate",
			Inner:   true,
			Handler: c.Secure.Authenticate,
		},
		{
			Path:    "system/secure/authorize",
			Inner:   true,
			Handler: c.Secure.Authorize,
		},
	}
}

func accessListCluster(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{
		{
			Path:    "system/access_list/get_by_id",
			Inner:   true,
			Handler: c.AccessList.GetById,
		},
		{
			Path:    "system/access_list/set_one",
			Inner:   true,
			Handler: c.AccessList.SetOne,
		},
		{
			Path:    "system/access_list/set_list",
			Inner:   true,
			Handler: c.AccessList.SetList,
		},
		{
			Path:    "system/access_list/delete_list",
			Inner:   true,
			Handler: c.AccessList.DeleteList,
		},
	}
}

func domainCluster(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{
		{
			Path:    "system/domain/get_domains_by_system_id",
			Inner:   true,
			Handler: c.Domain.GetBySystemId,
		},
		{
			Path:    "system/domain/create_update_domain",
			Inner:   true,
			Handler: c.Domain.CreateUpdate,
		},
		{
			Path:    "system/domain/get_domain_by_id",
			Inner:   true,
			Handler: c.Domain.GetById,
		},
		{
			Path:    "system/domain/delete_domains",
			Inner:   true,
			Handler: c.Domain.Delete,
		},
	}
}

// deprecated
func serviceCluster(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{
		{
			Path:    "system/service/get_service",
			Inner:   true,
			Handler: c.Service.Get,
		},
		{
			Path:    "system/service/get_services_by_domain_id",
			Inner:   true,
			Handler: c.Service.GetByDomainId,
		},
		{
			Path:    "system/service/create_update_service",
			Inner:   true,
			Handler: c.Service.CreateUpdate,
		},
		{
			Path:    "system/service/get_service_by_id",
			Inner:   true,
			Handler: c.Service.GetById,
		},
		{
			Path:    "system/service/delete_service",
			Inner:   true,
			Handler: c.Service.Delete,
		},
	}
}

func applicationCluster(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{
		{
			Path:    "system/application/get_applications",
			Inner:   true,
			Handler: c.Application.GetByIdList,
		},
		{
			Path:    "system/application/get_applications_by_service_id",
			Inner:   true,
			Handler: c.Application.GetByServiceId,
		},
		{
			Path:    "system/application/create_update_application",
			Inner:   true,
			Handler: c.Application.CreateUpdate,
		},
		{
			Path:    "system/application/get_application_by_id",
			Inner:   true,
			Handler: c.Application.GetById,
		},
		{
			Path:    "system/application/get_application_by_token",
			Inner:   true,
			Handler: c.Application.GetByToken,
		},
		{
			Path:    "system/application/delete_applications",
			Inner:   true,
			Handler: c.Application.Delete,
		},
		{
			Path:    "system/application/get_system_tree",
			Inner:   true,
			Handler: c.Application.GetSystemTree,
		},
		{
			Path:    "system/application/next_id",
			Inner:   true,
			Handler: c.Application.NextId,
		},
		{
			Path:    "system/application/get_all",
			Inner:   true,
			Handler: c.Application.GetAll,
		},
		{
			Path:    "system/application/create_application",
			Inner:   true,
			Handler: c.Application.Create,
		}, {
			Path:    "system/application/update_application",
			Inner:   true,
			Handler: c.Application.Update,
		},
	}
}

func tokenCluster(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{
		{
			Path:    "system/token/create_token",
			Inner:   true,
			Handler: c.Token.Create,
		},
		{
			Path:    "system/token/revoke_tokens",
			Inner:   true,
			Handler: c.Token.Revoke,
		},
		{
			Path:    "system/token/revoke_tokens_for_app",
			Inner:   true,
			Handler: c.Token.RevokeForApp,
		},
		{
			Path:    "system/token/get_tokens_by_app_id",
			Inner:   true,
			Handler: c.Token.GetByAppId,
		},
	}
}

func applicationGroupCluster(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{
		{
			Path:    "system/application_group/create",
			Inner:   true,
			Handler: c.AppGroup.Create,
		}, {
			Path:    "system/application_group/update",
			Inner:   true,
			Handler: c.AppGroup.Update,
		}, {
			Path:    "system/application_group/delete_list",
			Inner:   true,
			Handler: c.AppGroup.DeleteList,
		}, {
			Path:    "system/application_group/get_by_id_list",
			Inner:   true,
			Handler: c.AppGroup.GetByIdList,
		}, {
			Path:    "system/application_group/get_all",
			Inner:   true,
			Handler: c.AppGroup.GetAll,
		},
	}
}

func concatCluster(clusters ...[]cluster.EndpointDescriptor) []cluster.EndpointDescriptor {
	var result []cluster.EndpointDescriptor
	for _, c := range clusters {
		result = append(result, c...)
	}
	return result
}
