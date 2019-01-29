package model

import (
	"isp-system-service/entity"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	emptyService = &entity.Service{}
	ServiceRep   ServiceRepository
)

type ServiceRepository struct {
	DB orm.DB
}

func (sr *ServiceRepository) GetServices(list []int32) ([]entity.Service, error) {
	var res []entity.Service
	q := sr.DB.Model(&res)
	if len(res) > 0 {
		q = q.Where("id IN (?)", pg.In(list))
	}
	err := q.Order("created_at DESC").Select()
	return res, err
}

func (sr *ServiceRepository) GetServicesByDomainId(domainId ...int32) ([]entity.Service, error) {
	var res []entity.Service
	err := sr.DB.Model(&res).Where("domain_id IN (?)", pg.In(domainId)).Order("created_at DESC").Select()
	return res, err
}

func (sr *ServiceRepository) CreateService(service entity.Service) (entity.Service, error) {
	_, err := sr.DB.Model(&service).Insert()
	return service, err
}

func (sr *ServiceRepository) GetServiceByNameAndDomainId(name string, domainId int32) (*entity.Service, error) {
	service := new(entity.Service)
	err := sr.DB.Model(service).Where("name = ? AND domain_id = ?", name, domainId).First()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return service, err
}

func (sr *ServiceRepository) UpdateService(service entity.Service) (entity.Service, error) {
	_, err := sr.DB.Model(&service).Column("name", "description").WherePK().Returning("*").Update()
	return service, err
}

func (sr *ServiceRepository) GetServiceById(id int32) (*entity.Service, error) {
	service := &entity.Service{Id: id}
	err := sr.DB.Select(service)
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return service, err
}

func (sr *ServiceRepository) DeleteServices(list []int32) (int, error) {
	res, err := sr.DB.Model(emptyService).Where("id IN (?)", pg.In(list)).Delete()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
