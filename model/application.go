package model

import (
	"isp-system-service/entity"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	emptyApplication = &entity.Application{}
	AppRep           AppRepository
)

type AppRepository struct {
	DB orm.DB
}

func (ar *AppRepository) GetApplications(list []int32) ([]entity.Application, error) {
	res := make([]entity.Application, 0)
	q := ar.DB.Model(&res)
	if len(res) > 0 {
		q = q.Where("id IN (?)", pg.In(list))
	}
	err := q.Order("created_at DESC").Select()
	return res, err
}

func (ar *AppRepository) GetApplicationsByServiceId(serviceId ...int32) ([]entity.Application, error) {
	res := make([]entity.Application, 0)
	err := ar.DB.Model(&res).Where("service_id IN (?)", pg.In(serviceId)).Order("created_at DESC").Select()
	return res, err
}

func (ar *AppRepository) CreateApplication(service entity.Application) (entity.Application, error) {
	_, err := ar.DB.Model(&service).Insert()
	return service, err
}

func (ar *AppRepository) GetApplicationByNameAndServiceId(name string, serviceId int32) (*entity.Application, error) {
	app := new(entity.Application)
	err := ar.DB.Model(app).Where("name = ? AND service_id = ?", name, serviceId).First()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return app, err
}

func (ar *AppRepository) UpdateApplication(app entity.Application) (entity.Application, error) {
	_, err := ar.DB.Model(&app).Column("name", "description").WherePK().Returning("*").Update()
	return app, err
}

func (ar *AppRepository) GetApplicationById(id int32) (*entity.Application, error) {
	app := &entity.Application{Id: id}
	err := ar.DB.Select(app)
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return app, err
}

func (ar *AppRepository) DeleteApplications(list []int32) (int, error) {
	res, err := ar.DB.Model(emptyApplication).Where("id IN (?)", pg.In(list)).Delete()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
