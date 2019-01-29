package controller

import (
	"isp-system-service/entity"
	"isp-system-service/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetService(list []int32) ([]entity.Service, error) {
	res, err := model.ServiceRep.GetServices(list)
	if err != nil {
		return res, err
	}
	return res, nil
}

func GetServicesByDomainId(identity Identity) ([]entity.Service, error) {
	return model.ServiceRep.GetServicesByDomainId(identity.Id)
}

func CreateUpdateService(service entity.Service) (*entity.Service, error) {
	existed, err := model.ServiceRep.GetServiceByNameAndDomainId(service.Name, service.DomainId)
	if err != nil {
		return nil, err
	}
	domain, e := model.DomainRep.GetDomainById(service.DomainId)
	if e != nil {
		return nil, err
	}
	if domain == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Domain with id %d not found", service.DomainId)
	}
	if service.Id == 0 {
		if existed != nil {
			return nil, status.Errorf(codes.AlreadyExists, "Service with name %s already exists", service.Name)
		}
		service, e := model.ServiceRep.CreateService(service)
		return &service, e
	} else {
		if existed != nil && existed.Id != service.Id {
			return nil, status.Errorf(codes.AlreadyExists, "Service with name %s already exists", service.Name)
		}
		existed, err = model.ServiceRep.GetServiceById(service.Id)
		if err != nil {
			return nil, err
		}
		if existed == nil {
			return nil, status.Errorf(codes.NotFound, "Service with id %d not found", service.Id)
		}
		service, e := model.ServiceRep.UpdateService(service)
		return &service, e
	}
}

func GetServiceById(identity Identity) (*entity.Service, error) {
	service, err := model.ServiceRep.GetServiceById(identity.Id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, status.Errorf(codes.NotFound, "Service with id %d not found", identity.Id)
	}
	return service, err
}

func DeleteService(list []int32) (DeleteResponse, error) {
	if len(list) == 0 {
		return DeleteResponse{Deleted: 0}, status.Error(codes.InvalidArgument, "At least one id are required")
	}
	res, err := model.ServiceRep.DeleteServices(list)
	return DeleteResponse{Deleted: res}, err
}
