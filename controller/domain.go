package controller

import (
	"isp-system-service/entity"
	"isp-system-service/model"

	"github.com/integration-system/isp-lib/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

/*func GetDomains(list []int32) ([]entity.Domain, error) {
	res, err := model.DomainRep.GetDomains(list)
	if err != nil {
		return res, err
	}
	return res, nil
}*/

func GetDomainsBySystemId(md metadata.MD) ([]entity.Domain, error) {
	sysId, err := utils.ResolveMetadataIdentity(utils.SystemIdHeader, md)
	if err != nil {
		return nil, err
	}
	return model.DomainRep.GetDomainsBySystemId(int32(sysId))
}

func CreateUpdateDomain(domain entity.Domain, md metadata.MD) (*entity.Domain, error) {
	existed, err := model.DomainRep.GetDomainByNameAndSystemId(domain.Name, domain.SystemId)
	if err != nil {
		return nil, err
	}

	sysId, err := utils.ResolveMetadataIdentity(utils.SystemIdHeader, md)
	if err != nil {
		return nil, err
	}
	domain.SystemId = int32(sysId)

	sys, e := model.SystemRep.GetSystemById(domain.SystemId)
	if e != nil {
		return nil, err
	}
	if sys == nil {
		return nil, status.Errorf(codes.InvalidArgument, "System with id %d not found", domain.SystemId)
	}
	if domain.Id == 0 {
		if existed != nil {
			return nil, status.Errorf(codes.AlreadyExists, "Domain with name %s already exists", domain.Name)
		}
		domain, e := model.DomainRep.CreateDomain(domain)
		return &domain, e
	} else {
		if existed != nil && existed.Id != domain.Id {
			return nil, status.Errorf(codes.AlreadyExists, "Domain with name %s already exists", domain.Name)
		}
		existed, err = model.DomainRep.GetDomainById(domain.Id)
		if err != nil {
			return nil, err
		}
		if existed == nil {
			return nil, status.Errorf(codes.NotFound, "Domain with id %d not found", domain.Id)
		}
		domain, e := model.DomainRep.UpdateDomain(domain)
		return &domain, e
	}
}

func GetDomainById(identity Identity) (*entity.Domain, error) {
	domain, err := model.DomainRep.GetDomainById(identity.Id)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, status.Errorf(codes.NotFound, "Domain with id %d not found", identity.Id)
	}
	return domain, err
}

func DeleteDomains(list []int32) (DeleteResponse, error) {
	if len(list) == 0 {
		return DeleteResponse{Deleted: 0}, status.Errorf(codes.InvalidArgument, "At least one id are required")
	}
	res, err := model.DomainRep.DeleteDomains(list)
	return DeleteResponse{Deleted: res}, err
}
