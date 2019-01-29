package model

import (
	"isp-system-service/entity"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	emptyDomain = &entity.Domain{}
	DomainRep   DomainRepository
)

type DomainRepository struct {
	DB orm.DB
}

func (dr *DomainRepository) GetDomains(list []int32) ([]entity.Domain, error) {
	var res []entity.Domain
	q := dr.DB.Model(&res)
	if len(res) > 0 {
		q = q.Where("id IN (?)", pg.In(list))
	}
	err := q.Order("created_at DESC").Select()
	return res, err
}

func (dr *DomainRepository) GetDomainsBySystemId(systemId int32) ([]entity.Domain, error) {
	var res []entity.Domain
	err := dr.DB.Model(&res).Where("system_id = ?", systemId).Order("created_at DESC").Select()
	return res, err
}

func (dr *DomainRepository) CreateDomain(domain entity.Domain) (entity.Domain, error) {
	_, err := dr.DB.Model(&domain).Insert()
	return domain, err
}

func (dr *DomainRepository) GetDomainByNameAndSystemId(name string, systemId int32) (*entity.Domain, error) {
	domain := new(entity.Domain)
	err := dr.DB.Model(domain).Where("name = ? AND system_id = ?", name, systemId).First()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return domain, err
}

func (dr *DomainRepository) UpdateDomain(domain entity.Domain) (entity.Domain, error) {
	_, err := dr.DB.Model(&domain).Column("name", "description").WherePK().Returning("*").Update()
	return domain, err
}

func (dr *DomainRepository) GetDomainById(id int32) (*entity.Domain, error) {
	domain := &entity.Domain{Id: id}
	err := dr.DB.Select(domain)
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return domain, err
}

func (dr *DomainRepository) DeleteDomains(list []int32) (int, error) {
	res, err := dr.DB.Model(emptyDomain).Where("id IN (?)", pg.In(list)).Delete()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
