package migrations

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pressly/goose"
	"isp-system-service/domain"
)

const (
	insertDomain           = "INSERT INTO %s.domain (id, name, description, system_id) VALUES (%d, '%s', '%s', %d) RETURNING id;"
	insertApplicationGroup = "INSERT INTO %s.application_group (id, name, description, domain_id) VALUES (%d, '%s', '%s', %d) RETURNING id;"
	insertApp              = "INSERT INTO %s.application (id, name, description, application_group_id, type) VALUES (%d, '%s', '%s', %d, '%s') RETURNING id;"
	insertToken            = "INSERT INTO %s.token (token, app_id, expire_time) VALUES ('%s', %d, %d);"

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
	domainId           int
	appId              int
	applicationGroupId int
}

func (i *initializeMigration) SetParams(migrationDir string, schema string) {
	i.migrationDir = migrationDir
	i.schema = schema
}

func (i *initializeMigration) Up(tx *sql.Tx) error {
	path := filepath.Join(i.migrationDir, initFile)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return errors.WithMessage(err, "read file")
	}

	list := make([]domain.DomainWithApplicationGroup, 0)
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return errors.WithMessage(err, "json unmarshal")
	}

	lastDomainId := 0
	lastApplicationGroupId := 0
	lastAppId := 0
	for _, domainInfo := range list {
		last, err := i.insert(tx, domainInfo)
		if err != nil {
			return errors.WithMessage(err, "")
		}
		if last.domainId > lastDomainId {
			lastDomainId = last.domainId
		}
		if last.applicationGroupId > lastApplicationGroupId {
			lastApplicationGroupId = last.applicationGroupId
		}
		if last.appId > lastAppId {
			lastAppId = last.appId
		}
	}

	err = i.setSeq(tx, "domain", lastDomainId)
	if err != nil {
		return errors.WithMessage(err, "set domain")
	}

	err = i.setSeq(tx, "application_group", lastApplicationGroupId)
	if err != nil {
		return errors.WithMessage(err, "set application group")
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

func (i *initializeMigration) insert(tx *sql.Tx, domainInfo domain.DomainWithApplicationGroup) (*identity, error) {
	domainId, err := i.insertNode(tx, insertDomain, true, i.schema, domainInfo.Id, domainInfo.Name, domainInfo.Description, 1)
	if err != nil {
		return nil, errors.WithMessage(err, "insert domain")
	}

	applicationGroupId := 0
	appId := 0
	for _, applicationGroup := range domainInfo.ApplicationGroup {
		applicationGroupId, err = i.insertNode(tx, insertApplicationGroup, true, i.schema, applicationGroup.Id, applicationGroup.Name, applicationGroup.Description, domainId)
		if err != nil {
			return nil, errors.WithMessage(err, "insert application group")
		}

		for _, app := range applicationGroup.Apps {
			appId, err = i.insertNode(tx, insertApp, true, i.schema, app.Id, app.Name, app.Description, applicationGroupId, app.Type)
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
		domainId:           domainId,
		appId:              appId,
		applicationGroupId: applicationGroupId,
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
