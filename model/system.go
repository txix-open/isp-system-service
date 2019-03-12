package model

import (
	"isp-system-service/entity"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/integration-system/isp-lib/database"
)

var (
	emptySystem = &entity.System{}
	SystemRep   SystemRepository
)

type SystemRepository struct {
	DB orm.DB
}

func (sr *SystemRepository) GetSystems(list []int32) ([]entity.System, error) {
	res := make([]entity.System, 0)
	q := sr.DB.Model(&res)
	if len(res) > 0 {
		q = q.Where("id IN (?)", pg.In(list))
	}
	err := q.Order("created_at DESC").Select()
	return res, err
}

func (sr *SystemRepository) CreateSystem(system entity.System) (entity.System, error) {
	_, err := sr.DB.Model(&system).Insert()
	return system, err
}

func (sr *SystemRepository) GetSystemByName(name string) (*entity.System, error) {
	sys := new(entity.System)
	err := sr.DB.Model(sys).Where("name = ?", name).First()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return sys, err
}

func (sr *SystemRepository) UpdateSystem(system entity.System) (entity.System, error) {
	_, err := sr.DB.Model(&system).Column("name", "description").WherePK().Returning("*").Update()
	return system, err
}

func (sr *SystemRepository) GetSystemById(id int32) (*entity.System, error) {
	sys := &entity.System{Id: id}
	err := sr.DB.Select(sys)
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return sys, err
}

func (sr *SystemRepository) DeleteSystems(list []int32) (int, error) {
	res, err := sr.DB.Model(emptySystem).Where("id IN (?)", pg.In(list)).Delete()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func InitDbManager(pdb *database.DBManager) {
	db := pdb.Db
	SystemRep = SystemRepository{db}
	DomainRep = DomainRepository{db}
	ServiceRep = ServiceRepository{db}
	AppRep = AppRepository{db}
	TokenRep = TokenRepository{db}
}
