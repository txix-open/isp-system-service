package controller

import (
	"fmt"

	"isp-system-service/entity"
	"isp-system-service/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetSystems(list []int32) ([]entity.System, error) {
	res, err := model.SystemRep.GetSystems(list)
	if err != nil {
		return res, err
	}
	return res, nil
}

func CreateUpdateSystem(system entity.System) (*entity.System, error) {
	existed, err := model.SystemRep.GetSystemByName(system.Name)
	if err != nil {
		return nil, err
	}
	if system.Id == 0 {
		if existed != nil {
			return nil, status.Errorf(codes.AlreadyExists, "System with name %s already exists", system.Name)
		}
		sys, e := model.SystemRep.CreateSystem(system)
		return &sys, e
	} else {
		if existed != nil && existed.Id != system.Id {
			return nil, status.Errorf(codes.AlreadyExists, "System with name %s already exists", system.Name)
		}
		existed, err = model.SystemRep.GetSystemById(system.Id)
		if err != nil {
			return nil, err
		}
		if existed == nil {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("System with id %d not found", system.Id))
		}
		sys, e := model.SystemRep.UpdateSystem(system)
		return &sys, e
	}
}

func GetSystemById(identity Identity) (*entity.System, error) {
	sys, err := model.SystemRep.GetSystemById(identity.Id)
	if err != nil {
		return nil, err
	}
	if sys == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("System with id %d not found", identity.Id))
	}
	return sys, err
}

func DeleteSystems(list []int32) (DeleteResponse, error) {
	if len(list) == 0 {
		return DeleteResponse{}, status.Error(codes.InvalidArgument, "At least one id are required")
	}
	res, err := model.SystemRep.DeleteSystems(list)
	return DeleteResponse{Deleted: res}, err
}
