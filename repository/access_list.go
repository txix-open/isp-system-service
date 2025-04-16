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

type AccessList struct {
	db db.DB
}

func NewAccessList(db db.DB) AccessList {
	return AccessList{
		db: db,
	}
}

func (r AccessList) GetAccessListByAppIdAndMethod(ctx context.Context, appId int, method string) (*entity.AccessList, error) {
	q := `
	SELECT app_id, method, value
	FROM access_list
	WHERE app_id = $1 AND method = $2
`
	result := entity.AccessList{}
	err := r.db.SelectRow(ctx, &result, q, appId, method)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrAccessListNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r AccessList) GetAccessListByAppId(ctx context.Context, appId int) ([]entity.AccessList, error) {
	q := `
	SELECT app_id, method, value
	FROM access_list
	WHERE app_id = $1
`
	result := make([]entity.AccessList, 0)
	err := r.db.Select(ctx, &result, q, appId)
	if err != nil {
		return nil, errors.WithMessagef(err, "select db")
	}

	return result, nil
}

func (r AccessList) GetAccessListByAppIdList(ctx context.Context, appIdList []int) ([]entity.AccessList, error) {
	q, args, err := query.New().
		Select("app_id", "method", "value").
		From("access_list").
		Where(squirrel.Eq{"app_id": appIdList}).
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.AccessList, 0)
	err = r.db.Select(ctx, &result, q, args...)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}

func (r AccessList) InsertArrayAccessList(ctx context.Context, entity []entity.AccessList) error {
	qBuilder := query.New().
		Insert("access_list").
		Columns("app_id", "method", "value")
	for _, e := range entity {
		qBuilder = qBuilder.Values(e.AppId, e.Method, e.Value)
	}
	q, args, err := qBuilder.ToSql()
	if err != nil {
		return errors.WithMessage(err, "build query")
	}

	_, err = r.db.Exec(ctx, q, args...)
	if err != nil {
		return errors.WithMessagef(err, "exec query %s", q)
	}

	return nil
}

func (r AccessList) UpsertAccessList(ctx context.Context, e entity.AccessList) (int, error) {
	q := `
	INSERT INTO access_list 
	(app_id, method, value)
	VALUES ($1, $2, $3)
	ON CONFLICT (app_id, method) 
		DO UPDATE
			SET (value) = (SELECT EXCLUDED.value)
`
	result, err := r.db.Exec(ctx, q, e.AppId, e.Method, e.Value)
	if err != nil {
		return 0, errors.WithMessagef(err, "exec query %s", q)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.WithMessage(err, "get rows affected")
	}

	return int(rowsAffected), nil
}

func (r AccessList) DeleteAccessList(ctx context.Context, appId int, methods []string) error {
	q, args, err := query.New().
		Delete("access_list").
		Where(squirrel.Eq{
			"app_id": appId,
			"method": methods,
		}).
		Suffix("RETURNING app_id, method, value").
		ToSql()
	if err != nil {
		return errors.WithMessage(err, "build query")
	}

	_, err = r.db.Exec(ctx, q, args...)
	if err != nil {
		return errors.WithMessagef(err, "exec db")
	}

	return nil
}

func (r AccessList) DeleteAccessListByAppId(ctx context.Context, appId int) ([]entity.AccessList, error) {
	q := `
	DELETE FROM access_list
	WHERE app_id = $1
	RETURNING app_id, method, value
`
	result := make([]entity.AccessList, 0)
	err := r.db.Select(ctx, &result, q, appId)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}
