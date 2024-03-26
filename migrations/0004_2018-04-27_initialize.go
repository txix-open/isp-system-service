package migrations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pressly/goose"
	"isp-system-service/domain"
)

const (
	insertDomain  = "INSERT INTO %s.domain (id, name, description, system_id) VALUES (%d, '%s', '%s', %d) RETURNING id;"
	insertService = "INSERT INTO %s.service (id, name, description, domain_id) VALUES (%d, '%s', '%s', %d) RETURNING id;"
	insertApp     = "INSERT INTO %s.application (id, name, description, service_id, type) VALUES (%d, '%s', '%s', %d, '%s') RETURNING id;"
	insertToken   = "INSERT INTO %s.token (token, app_id, expire_time) VALUES ('%s', %d, %d);"

	setSeqQuery = "SELECT setval('%s.%s_id_seq' :: regclass, %d);"

	initFile = "init.json"
)

func init() {
	goose.AddMigration(Initialize.Up, Initialize.Down)
}

var Initialize = &initializeMigration{}

type initializeMigration struct {
	migrationDir string
	schema       string
}

type identity struct {
	domainId  int
	appId     int
	serviceId int
}

func (i *initializeMigration) SetParams(migrationDir string, schema string) {
	i.migrationDir = migrationDir
	i.schema = schema
}

func (i *initializeMigration) Up(tx *sql.Tx) error {
	path := filepath.Join(i.migrationDir, initFile)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.WithMessage(err, "read file")
	}

	list := make([]domain.DomainWithService, 0)
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return errors.WithMessage(err, "json unmarshal")
	}

	lastDomainId := 0
	lastServiceId := 0
	lastAppId := 0
	for _, domainInfo := range list {
		last, err := i.insert(tx, domainInfo)
		if err != nil {
			return errors.WithMessage(err, "")
		}
		if last.domainId > lastDomainId {
			lastDomainId = last.domainId
		}
		if last.serviceId > lastServiceId {
			lastServiceId = last.serviceId
		}
		if last.appId > lastAppId {
			lastAppId = last.appId
		}
	}

	err = i.setSeq(tx, "domain", lastDomainId)
	if err != nil {
		return errors.WithMessage(err, "set domain")
	}

	err = i.setSeq(tx, "service", lastServiceId)
	if err != nil {
		return errors.WithMessage(err, "set service")
	}

	err = i.setSeq(tx, "application", lastAppId)
	if err != nil {
		return errors.WithMessage(err, "set application")
	}

	return nil
}

func (i *initializeMigration) Down(tx *sql.Tx) error {
	return nil
}

func (i *initializeMigration) insert(tx *sql.Tx, domainInfo domain.DomainWithService) (*identity, error) {
	domainId, err := i.insertNode(tx, insertDomain, true, i.schema, domainInfo.Id, domainInfo.Name, domainInfo.Description, 1)
	if err != nil {
		return nil, errors.WithMessage(err, "insert domain")
	}

	serviceId := 0
	appId := 0
	for _, service := range domainInfo.Services {
		serviceId, err = i.insertNode(tx, insertService, true, i.schema, service.Id, service.Name, service.Description, domainId)
		if err != nil {
			return nil, errors.WithMessage(err, "insert service")
		}

		for _, app := range service.Apps {
			appId, err = i.insertNode(tx, insertApp, true, i.schema, app.Id, app.Name, app.Description, serviceId, app.Type)
			if err != nil {
				return nil, errors.WithMessage(err, "insert application")
			}

			for _, t := range app.Tokens {
				_, err = i.insertNode(tx, insertToken, false, i.schema, t.Token, appId, t.ExpireTime)
				if err != nil {
					return nil, errors.WithMessage(err, "insert token")
				}
			}
		}
	}

	return &identity{
		domainId:  domainId,
		appId:     appId,
		serviceId: serviceId,
	}, nil
}

func (i *initializeMigration) insertNode(tx *sql.Tx, query string, scan bool, args ...interface{}) (int, error) {
	var id int
	query = fmt.Sprintf(query, args...)
	row := tx.QueryRow(query)
	if scan {
		err := row.Scan(&id)
		if err != nil {
			return 0, errors.WithMessage(err, "row scan")
		}
	}

	return id, nil
}

func (i *initializeMigration) setSeq(tx *sql.Tx, table string, value int) error {
	q := fmt.Sprintf(setSeqQuery, i.schema, table, value)
	_, err := tx.Exec(q)
	if err != nil {
		return errors.WithMessage(err, "exec db")
	}

	return nil
}
