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

type AppGroup struct {
	db db.DB
}

func NewAppGroup(db db.DB) AppGroup {
	return AppGroup{
		db: db,
	}
}

func (r AppGroup) GetAppGroupById(ctx context.Context, id int) (*entity.AppGroup, error) {
	q := `
	SELECT id, name, description, domain_id, created_at, updated_at
	FROM application_group
	WHERE id = $1
`
	result := entity.AppGroup{}
	err := r.db.SelectRow(ctx, &result, q, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrAppGroupNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r AppGroup) GetAppGroupByIdList(ctx context.Context, idList []int) ([]entity.AppGroup, error) {
	q, arg, err := query.New().
		Select("id", "name", "description", "domain_id", "created_at", "updated_at").
		From("application_group").
		Where(squirrel.Eq{"id": idList}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.AppGroup, 0)
	err = r.db.Select(ctx, &result, q, arg...)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}

func (r AppGroup) GetAppGroupByDomainId(ctx context.Context, domainIdList []int) ([]entity.AppGroup, error) {
	q, arg, err := query.New().
		Select("id", "name", "description", "domain_id", "created_at", "updated_at").
		From("application_group").
		Where(squirrel.Eq{"domain_id": domainIdList}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.AppGroup, 0)
	err = r.db.Select(ctx, &result, q, arg...)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}

func (r AppGroup) GetAppGroupByNameAndDomainId(ctx context.Context, name string, domainId int) (*entity.AppGroup, error) {
	q := `
	SELECT id, name, description, domain_id, created_at, updated_at
	FROM application_group
	WHERE name = $1 AND domain_id = $2
`
	result := entity.AppGroup{}
	err := r.db.SelectRow(ctx, &result, q, name, domainId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrAppGroupNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r AppGroup) CreateAppGroup(ctx context.Context, name string, desc string, domainId int) (*entity.AppGroup, error) {
	q := `
	INSERT INTO application_group
	(name, description, domain_id)
	VALUES ($1, $2, $3)
	ON CONFLICT (name, domain_id) DO NOTHING
	RETURNING id, name, description, domain_id, created_at, updated_at
`
	result := entity.AppGroup{}
	err := r.db.SelectRow(ctx, &result, q, name, desc, domainId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrAppGroupDuplicateName
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r AppGroup) UpdateAppGroup(ctx context.Context, id int, name string, description string) (*entity.AppGroup, error) {
	q := `
	UPDATE application_group 
	SET name = $1, description = $2
	WHERE id = $3
	RETURNING id, name, description, domain_id, created_at, updated_at
`
	result := entity.AppGroup{}
	err := r.db.SelectRow(ctx, &result, q, name, description, id)
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.ConstraintName == applicationGroupUniqueNameConstraint:
		return nil, domain.ErrAppGroupDuplicateName
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrAppGroupNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r AppGroup) DeleteAppGroup(ctx context.Context, idList []int) (int, error) {
	q, args, err := query.New().
		Delete("application_group").
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

func (r AppGroup) GetAllAppGroups(ctx context.Context) ([]entity.AppGroup, error) {
	q, arg, err := query.New().
		Select("id", "name", "description", "domain_id", "created_at", "updated_at").
		From("application_group").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.AppGroup, 0)
	err = r.db.Select(ctx, &result, q, arg...)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}
