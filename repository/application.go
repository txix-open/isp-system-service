package repository

import (
	"context"
	"database/sql"

	"isp-system-service/domain"
	"isp-system-service/entity"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/db"
	"github.com/txix-open/isp-kit/db/query"
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
	SELECT id, name, description, application_group_id, type, created_at, updated_at
	FROM application
	WHERE id = $1
`
	result := entity.Application{}
	err := r.db.SelectRow(ctx, &result, q, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrApplicationNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r Application) GetApplicationByIdList(ctx context.Context, idList []int) ([]entity.Application, error) {
	q, args, err := query.New().
		Select("id", "name", "description", "application_group_id", "type", "created_at", "updated_at").
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
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}

func (r Application) GetApplicationByAppGroupIdList(ctx context.Context, appGroupIdList []int) ([]entity.Application, error) {
	q, args, err := query.New().
		Select("id", "name", "description", "application_group_id", "type", "created_at", "updated_at").
		From("application").
		Where(squirrel.Eq{"application_group_id": appGroupIdList}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	result := make([]entity.Application, 0)
	err = r.db.Select(ctx, &result, q, args...)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}

func (r Application) GetApplicationByNameAndAppGroupId(ctx context.Context, name string, serviceId int) (*entity.Application, error) {
	q := `
	SELECT id, name, description, application_group_id, type, created_at, updated_at
	FROM application 
	WHERE name = $1 AND application_group_id = $2
`
	result := entity.Application{}
	err := r.db.SelectRow(ctx, &result, q, name, serviceId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil // nolint:nilnil
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r Application) CreateApplication(ctx context.Context, id int, name string, desc string, appGroupId int, appType string) (
	*entity.Application, error) {
	q := `
	INSERT INTO application 
	(id, name, description, application_group_id, type)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, name, description, application_group_id, type, created_at, updated_at
`
	result := entity.Application{}
	err := r.db.SelectRow(ctx, &result, q, id, name, desc, appGroupId, appType)
	if err != nil {
		return nil, r.handleCreateError(err, q)
	}
	return &result, nil
}

func (r Application) UpdateApplication(
	ctx context.Context,
	id int,
	name string,
	description string,
) (*entity.Application, error) {
	q := `
	UPDATE application 
	SET name = $2,
		description = $3
	WHERE id = $1
	RETURNING id, name, description, application_group_id, type, created_at, updated_at
`
	result := entity.Application{}
	err := r.db.SelectRow(ctx, &result, q, id, name, description)
	if err != nil {
		return nil, r.handleUpdateError(err, q)
	}
	return &result, nil
}

func (r Application) UpdateApplicationWithNewId(
	ctx context.Context,
	oldId int,
	newId int,
	name string,
	description string,
) (*entity.Application, error) {
	q := `
	UPDATE application 
	SET id = $2,
		name = $3,
		description = $4
	WHERE id = $1
	RETURNING id, name, description, application_group_id, type, created_at, updated_at
`
	result := entity.Application{}
	err := r.db.SelectRow(ctx, &result, q, oldId, newId, name, description)
	if err != nil {
		return nil, r.handleUpdateError(err, q)
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

func (r Application) NextApplicationId(ctx context.Context) (int, error) {
	q := `SELECT COALESCE(MAX(id)+1, 1) FROM application;`
	nextAppId := 0
	err := r.db.SelectRow(ctx, &nextAppId, q)
	if err != nil {
		return -1, errors.WithMessagef(err, "exec query %s", q)
	}
	return nextAppId, nil
}

func (r Application) GetAllApplications(ctx context.Context) ([]entity.Application, error) {
	q, args, err := query.New().
		Select("id", "name", "description", "application_group_id", "type", "created_at", "updated_at").
		From("application").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	result := make([]entity.Application, 0)
	err = r.db.Select(ctx, &result, q, args...)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}

func (r Application) handleCreateError(err error, q string) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return errors.WithMessagef(err, "exec query %s", q)
	}
	switch pgErr.ConstraintName {
	case applicationPkConstraintName:
		return domain.ErrApplicationDuplicateId
	case applicationUniqueNameConstraintName:
		return domain.ErrApplicationDuplicateName
	case applicationFkAppGroupConstraintName:
		return domain.ErrAppGroupNotFound
	}
	return errors.WithMessagef(err, "exec query %s", q)
}

func (r Application) handleUpdateError(err error, q string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrApplicationNotFound
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return errors.WithMessagef(err, "exec query %s", q)
	}

	switch pgErr.ConstraintName {
	case applicationPkConstraintName:
		return domain.ErrApplicationDuplicateId
	case applicationUniqueNameConstraintName:
		return domain.ErrApplicationDuplicateName
	default:
		return errors.WithMessagef(err, "exec query %s", q)
	}
}
