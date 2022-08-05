package tests

import (
	"github.com/integration-system/isp-kit/test/dbt"
	"isp-system-service/entity"
)

func InsertDomain(db *dbt.TestDb, value entity.Domain) {
	q := `
	INSERT INTO domain 
		(id, name, description, system_id, created_at, updated_at)
	VALUES 
		(:id, :name, :description, 1, :created_at, :updated_at)
`
	db.Must().ExecNamed(q, value)
}

func InsertService(db *dbt.TestDb, value entity.Service) {
	q := `
	INSERT INTO service 
		(id, name, description, domain_id, created_at, updated_at)
	VALUES 
		(:id, :name, :description, :domain_id, :created_at, :updated_at)
`
	db.Must().ExecNamed(q, value)
}

func InsertApplication(db *dbt.TestDb, value entity.Application) {
	q := `
	INSERT INTO application 
		(id, name, description, service_id, created_at, updated_at) 
	VALUES 
		(:id, :name, :description, :service_id, :created_at, :updated_at)
`
	db.Must().ExecNamed(q, value)
}

func InsertToken(db *dbt.TestDb, value entity.Token) {
	q := `
	INSERT INTO token 
		(token, created_at, app_id, expire_time) 
	VALUES 
		(:token, :created_at, :app_id, :expire_time)
`
	db.Must().ExecNamed(q, value)
}

func InsertAccessList(db *dbt.TestDb, value entity.AccessList) {
	q := `
	INSERT INTO access_list
	    (app_id, method, value)
	VALUES
		(:app_id, :method, :value)
`
	db.Must().ExecNamed(q, value)
}
