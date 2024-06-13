package repository

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/db"
	"github.com/txix-open/isp-kit/db/query"
	"isp-system-service/domain"
	"isp-system-service/entity"
)

type ApplicationGroup struct {
	db db.DB
}

func NewApplicationGroup(db db.DB) ApplicationGroup {
	return ApplicationGroup{
		db: db,
	}
}

func (r ApplicationGroup) GetApplicationGroupById(ctx context.Context, id int) (*entity.ApplicationGroup, error) {
	q := `
	SELECT id, name, description, domain_id, created_at, updated_at
	FROM application_group
	WHERE id = $1
`
	result := entity.ApplicationGroup{}
	err := r.db.SelectRow(ctx, &result, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrApplicationGroupNotFound
		}
		return nil, errors.WithMessage(err, "select row db")
	}

	return &result, nil
}

func (r ApplicationGroup) GetApplicationGroupByIdList(ctx context.Context, idList []int) ([]entity.ApplicationGroup, error) {
	q, arg, err := query.New().
		Select("id", "name", "description", "domain_id", "created_at", "updated_at").
		From("application_group").
		Where(squirrel.Eq{"id": idList}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.ApplicationGroup, 0)
	err = r.db.Select(ctx, &result, q, arg...)
	if err != nil {
		return nil, errors.WithMessagef(err, "select db")
	}

	return result, nil
}

func (r ApplicationGroup) GetApplicationGroupByDomainId(ctx context.Context, domainIdList []int) ([]entity.ApplicationGroup, error) {
	q, arg, err := query.New().
		Select("id", "name", "description", "domain_id", "created_at", "updated_at").
		From("application_group").
		Where(squirrel.Eq{"domain_id": domainIdList}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.ApplicationGroup, 0)
	err = r.db.Select(ctx, &result, q, arg...)
	if err != nil {
		return nil, errors.WithMessagef(err, "select db")
	}

	return result, nil
}

func (r ApplicationGroup) GetSApplicationGroupByNameAndDomainId(ctx context.Context, name string, domainId int) (*entity.ApplicationGroup, error) {
	q := `
	SELECT id, name, description, domain_id, created_at, updated_at
	FROM application_group
	WHERE name = $1 AND domain_id = $2
`
	result := entity.ApplicationGroup{}
	err := r.db.SelectRow(ctx, &result, q, name, domainId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrApplicationGroupNotFound
		}
		return nil, errors.WithMessage(err, "select db")
	}

	return &result, nil
}

func (r ApplicationGroup) CreateApplicationGroup(ctx context.Context, name string, desc string, domainId int) (*entity.ApplicationGroup, error) {
	q := `
	INSERT INTO application_group
	(name, description, domain_id)
	VALUES ($1, $2, $3)
	RETURNING id, name, description, domain_id, created_at, updated_at
`
	result := entity.ApplicationGroup{}
	err := r.db.SelectRow(ctx, &result, q, name, desc, domainId)
	if err != nil {
		return nil, errors.WithMessage(err, "select row db")
	}

	return &result, nil
}

func (r ApplicationGroup) UpdateApplicationGroup(ctx context.Context, id int, name string, description string) (*entity.ApplicationGroup, error) {
	q := `
	UPDATE application_group 
	SET name = $1, description = $2
	WHERE id = $3
	RETURNING id, name, description, domain_id, created_at, updated_at
`
	result := entity.ApplicationGroup{}
	err := r.db.SelectRow(ctx, &result, q, name, description, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrApplicationGroupNotFound
		}
		return nil, errors.WithMessage(err, "select row db")
	}

	return &result, nil
}

func (r ApplicationGroup) DeleteApplicationGroup(ctx context.Context, idList []int) (int, error) {
	q, args, err := query.New().
		Delete("application_group").
		Where(squirrel.Eq{"id": idList}).
		ToSql()
	if err != nil {
		return 0, errors.WithMessagef(err, "build query")
	}

	result, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return 0, errors.WithMessagef(err, "exec db")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.WithMessagef(err, "get rows affected")
	}

	return int(rowsAffected), nil
}
