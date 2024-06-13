package routes

import (
	"github.com/txix-open/isp-kit/cluster"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/endpoint"
	"github.com/txix-open/isp-kit/grpc/isp"
	"isp-system-service/controller"
)

type Controllers struct {
	AccessList       controller.AccessList
	Domain           controller.Domain
	Service          controller.Service
	ApplicationGroup controller.ApplicationGroup
	Application      controller.Application
	Token            controller.Token
	Secure           controller.Secure
}

func EndpointDescriptors() []cluster.EndpointDescriptor {
	return endpointDescriptors(Controllers{})
}

func Handler(wrapper endpoint.Wrapper, c Controllers) isp.BackendServiceServer {
	muxer := grpc.NewMux()
	for _, descriptor := range endpointDescriptors(c) {
		muxer.Handle(descriptor.Path, wrapper.Endpoint(descriptor.Handler))
	}
	return muxer
}

func endpointDescriptors(c Controllers) []cluster.EndpointDescriptor {
	return concatCluster([][]cluster.EndpointDescriptor{
		secureCluster(c),
		accessListCluster(c),
		domainCluster(c),
		serviceCluster(c),
		applicationGroupCluster(c),
		applicationCluster(c),
		tokenCluster(c),
	})
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

func applicationGroupCluster(c Controllers) []cluster.EndpointDescriptor {
	return []cluster.EndpointDescriptor{
		{
			Path:    "system/application_group/get_group",
			Inner:   true,
			Handler: c.ApplicationGroup.Get,
		},
		{
			Path:    "system/application_group/create_update_group",
			Inner:   true,
			Handler: c.ApplicationGroup.CreateUpdate,
		},
		{
			Path:    "system/application_group/get_group_by_id",
			Inner:   true,
			Handler: c.ApplicationGroup.GetById,
		},
		{
			Path:    "system/application_group/delete_group",
			Inner:   true,
			Handler: c.ApplicationGroup.Delete,
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
			Path:    "system/application/get_applications_by_application_group_id",
			Inner:   true,
			Handler: c.Application.GetByApplicationGroupId,
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
			Path:    "system/application/delete_applications",
			Inner:   true,
			Handler: c.Application.Delete,
		},
		{
			Path:    "system/application/get_system_tree",
			Inner:   true,
			Handler: c.Application.GetSystemTree,
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

func concatCluster(clusters [][]cluster.EndpointDescriptor) []cluster.EndpointDescriptor {
	var result []cluster.EndpointDescriptor
	for _, c := range clusters {
		result = append(result, c...)
	}
	return result
}
