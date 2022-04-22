package routes

import (
	"github.com/integration-system/isp-kit/cluster"
	"github.com/integration-system/isp-kit/grpc"
	"github.com/integration-system/isp-kit/grpc/endpoint"
	"github.com/integration-system/isp-kit/grpc/isp"
	"isp-system-service/controller"
)

type Controllers struct {
	AccessList  controller.AccessList
	Domain      controller.Domain
	Service     controller.Service
	Application controller.Application
	Token       controller.Token
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
		accessListCluster(c),
		domainCluster(c),
		serviceCluster(c),
		applicationCluster(c),
		tokenCluster(c),
	})
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
