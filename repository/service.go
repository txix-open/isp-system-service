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

type Service struct {
	db db.DB
}

func NewService(db db.DB) Service {
	return Service{
		db: db,
	}
}

func (r Service) GetServiceById(ctx context.Context, id int) (*entity.Service, error) {
	q := `
	SELECT id, name, description, domain_id, created_at, updated_at
	FROM service
	WHERE id = $1
`
	result := entity.Service{}
	err := r.db.SelectRow(ctx, &result, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrServiceNotFound
		}
		return nil, errors.WithMessage(err, "select row db")
	}

	return &result, nil
}

func (r Service) GetServiceByIdList(ctx context.Context, idList []int) ([]entity.Service, error) {
	q, arg, err := query.New().
		Select("id", "name", "description", "domain_id", "created_at", "updated_at").
		From("service").
		Where(squirrel.Eq{"id": idList}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.Service, 0)
	err = r.db.Select(ctx, &result, q, arg...)
	if err != nil {
		return nil, errors.WithMessagef(err, "select db")
	}

	return result, nil
}

func (r Service) GetServiceByDomainId(ctx context.Context, domainIdList []int) ([]entity.Service, error) {
	q, arg, err := query.New().
		Select("id", "name", "description", "domain_id", "created_at", "updated_at").
		From("service").
		Where(squirrel.Eq{"domain_id": domainIdList}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.Service, 0)
	err = r.db.Select(ctx, &result, q, arg...)
	if err != nil {
		return nil, errors.WithMessagef(err, "select db")
	}

	return result, nil
}

func (r Service) GetServiceByNameAndDomainId(ctx context.Context, name string, domainId int) (*entity.Service, error) {
	q := `
	SELECT id, name, description, domain_id, created_at, updated_at
	FROM service
	WHERE name = $1 AND domain_id = $2
`
	result := entity.Service{}
	err := r.db.SelectRow(ctx, &result, q, name, domainId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrServiceNotFound
		}
		return nil, errors.WithMessage(err, "select db")
	}

	return &result, nil
}

func (r Service) CreateService(ctx context.Context, name string, desc string, domainId int) (*entity.Service, error) {
	q := `
	INSERT INTO service
	(name, description, domain_id)
	VALUES ($1, $2, $3)
	RETURNING id, name, description, domain_id, created_at, updated_at
`
	result := entity.Service{}
	err := r.db.SelectRow(ctx, &result, q, name, desc, domainId)
	if err != nil {
		return nil, errors.WithMessage(err, "select row db")
	}

	return &result, nil
}

func (r Service) UpdateService(ctx context.Context, id int, name string, description string) (*entity.Service, error) {
	q := `
	UPDATE service 
	SET name = $1, description = $2
	WHERE id = $3
	RETURNING id, name, description, domain_id, created_at, updated_at
`
	result := entity.Service{}
	err := r.db.SelectRow(ctx, &result, q, name, description, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrServiceNotFound
		}
		return nil, errors.WithMessage(err, "select row db")
	}

	return &result, nil
}

func (r Service) DeleteService(ctx context.Context, idList []int) (int, error) {
	q, args, err := query.New().
		Delete("service").
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
