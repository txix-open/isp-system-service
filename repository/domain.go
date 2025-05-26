package repository

import (
	"context"
	"database/sql"

	"isp-system-service/domain"
	"isp-system-service/entity"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/db"
	"github.com/txix-open/isp-kit/db/query"
)

type Domain struct {
	db db.DB
}

func NewDomain(db db.DB) Domain {
	return Domain{
		db: db,
	}
}

func (r Domain) GetDomainById(ctx context.Context, id int) (*entity.Domain, error) {
	q := `
	SELECT id, name, description, system_id, created_at, updated_at
	FROM domain
	WHERE id = $1
`
	result := entity.Domain{}
	err := r.db.SelectRow(ctx, &result, q, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrDomainNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r Domain) GetDomainByIdList(ctx context.Context, idList []int) ([]entity.Domain, error) {
	q, args, err := query.New().
		Select("id", "name", "description", "system_id", "created_at", "updated_at").
		From("domain").
		Where(squirrel.Eq{"id": idList}).
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.Domain, 0)
	err = r.db.Select(ctx, &result, q, args...)
	if err != nil {
		return nil, errors.WithMessage(err, "select db")
	}

	return result, nil
}

func (r Domain) GetDomainBySystemId(ctx context.Context, systemId int) ([]entity.Domain, error) {
	q := `
	SELECT id, name, description, system_id, created_at, updated_at
	FROM domain
	WHERE system_id = $1
	ORDER BY created_at DESC
`
	result := make([]entity.Domain, 0)
	err := r.db.Select(ctx, &result, q, systemId)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}

func (r Domain) GetDomainByNameAndSystemId(ctx context.Context, name string, systemId int) (*entity.Domain, error) {
	q := `
	SELECT id, name, description, system_id, created_at, updated_at
	FROM domain
	WHERE name = $1 AND system_id = $2 
`
	result := entity.Domain{}
	err := r.db.SelectRow(ctx, &result, q, name, systemId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrDomainNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r Domain) CreateDomain(ctx context.Context, name string, desc string, systemId int) (*entity.Domain, error) {
	q := `
	INSERT INTO domain
	(name, description, system_id)
	VALUES ($1, $2, $3)
	RETURNING id, name, description, system_id, created_at, updated_at
`
	result := entity.Domain{}
	err := r.db.SelectRow(ctx, &result, q, name, desc, systemId)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}
	return &result, nil
}

func (r Domain) UpdateDomain(ctx context.Context, id int, name string, description string) (*entity.Domain, error) {
	q := `
	UPDATE domain 
	SET name = $1, description = $2
	WHERE id = $3
	RETURNING id, name, description, system_id, created_at, updated_at
`
	result := entity.Domain{}
	err := r.db.SelectRow(ctx, &result, q, name, description, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrDomainNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r Domain) DeleteDomain(ctx context.Context, idList []int) (int, error) {
	q, args, err := query.New().
		Delete("domain").
		Where(squirrel.Eq{"id": idList}).
		ToSql()
	if err != nil {
		return 0, errors.WithMessagef(err, "build query")
	}

	result, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return 0, errors.WithMessagef(err, "exec query %s", q)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.WithMessagef(err, "get rows affected")
	}

	return int(rowsAffected), nil
}
