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

type Application struct {
	db db.DB
}

func NewApplication(db db.DB) Application {
	return Application{
		db: db,
	}
}

func (r Application) GetApplicationById(ctx context.Context, id int) (*entity.Application, error) {
	q := `
	SELECT id, name, description, service_id, type, created_at, updated_at
	FROM application
	WHERE id = $1
`
	result := entity.Application{}
	err := r.db.SelectRow(ctx, &result, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrApplicationNotFound
		}
		return nil, errors.WithMessage(err, "select row db")
	}

	return &result, nil
}

func (r Application) GetApplicationByIdList(ctx context.Context, idList []int) ([]entity.Application, error) {
	q, args, err := query.New().
		Select("id", "name", "description", "service_id", "type", "created_at", "updated_at").
		From("application").
		Where(squirrel.Eq{"id": idList}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.Application, 0)
	err = r.db.Select(ctx, &result, q, args...)
	if err != nil {
		return nil, errors.WithMessagef(err, "select db")
	}

	return result, nil
}

func (r Application) GetApplicationByServiceIdList(ctx context.Context, serviceIdList []int) ([]entity.Application, error) {
	q, args, err := query.New().
		Select("id", "name", "description", "service_id", "type", "created_at", "updated_at").
		From("application").
		Where(squirrel.Eq{"service_id": serviceIdList}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.Application, 0)
	err = r.db.Select(ctx, &result, q, args...)
	if err != nil {
		return nil, errors.WithMessagef(err, "select db")
	}

	return result, nil
}

func (r Application) GetApplicationByNameAndServiceId(ctx context.Context, name string, serviceId int) (*entity.Application, error) {
	q := `
	SELECT id, name, description, service_id, type, created_at, updated_at
	FROM application 
	WHERE name = $1 AND service_id = $2
`
	result := entity.Application{}
	err := r.db.SelectRow(ctx, &result, q, name, serviceId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}
		return nil, errors.WithMessagef(err, "select row db")
	}

	return &result, nil
}

func (r Application) CreateApplication(ctx context.Context, name string, desc string, serviceId int, appType string) (
	*entity.Application, error) {
	q := `
	INSERT INTO application 
	(name, description, service_id, type)
	VALUES ($1, $2, $3, $4)
	RETURNING id, name, description, service_id, type, created_at, updated_at
`
	result := entity.Application{}
	err := r.db.SelectRow(ctx, &result, q, name, desc, serviceId, appType)
	if err != nil {
		return nil, errors.WithMessagef(err, "select row db")
	}

	return &result, nil
}

func (r Application) UpdateApplication(ctx context.Context, id int, name string, description string) (*entity.Application, error) {
	q := `
	UPDATE application 
	SET name = $2,
		description = $3
	WHERE id = $1
	RETURNING id, name, description, service_id, type, created_at, updated_at
`
	result := entity.Application{}
	err := r.db.SelectRow(ctx, &result, q, id, name, description)
	if err != nil {
		return nil, errors.WithMessagef(err, "select row db")
	}

	return &result, nil
}

func (r Application) DeleteApplicationByIdList(ctx context.Context, idList []int) (int, error) {
	q, args, err := query.New().
		Delete("application").
		Where(squirrel.Eq{"id": idList}).
		ToSql()
	if err != nil {
		return 0, errors.WithMessage(err, "build query")
	}

	result, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return 0, errors.WithMessage(err, "exec db")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.WithMessage(err, "get rows affected")
	}

	return int(rowsAffected), nil
}
