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
	"github.com/txix-open/isp-kit/metrics/sql_metrics"
)

type AccessList struct {
	db db.DB
}

func NewAccessList(db db.DB) AccessList {
	return AccessList{
		db: db,
	}
}

func (r AccessList) GetAccessListByAppIdAndMethod(ctx context.Context, appId int, httpMethod string, method string) (*entity.AccessList, error) {
	ctx = sql_metrics.OperationLabelToContext(ctx, "AccessList.GetAccessListByAppIdAndMethod")

	q := `
		SELECT app_id, http_method, method, value
		FROM access_list
		WHERE app_id = $1
		AND method = $2
		AND http_method IN ($3, '')
		ORDER BY (http_method = $3) DESC
		LIMIT 1;
	`
	result := entity.AccessList{}
	err := r.db.SelectRow(ctx, &result, q, appId, method, httpMethod)
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
	ctx = sql_metrics.OperationLabelToContext(ctx, "AccessList.GetAccessListByAppId")

	q := `
	SELECT app_id, http_method, method, value
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
	ctx = sql_metrics.OperationLabelToContext(ctx, "AccessList.GetAccessListByAppIdList")

	q, args, err := query.New().
		Select("app_id", "http_method", "method", "value").
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
	ctx = sql_metrics.OperationLabelToContext(ctx, "AccessList.InsertArrayAccessList")

	qBuilder := query.New().
		Insert("access_list").
		Columns("app_id", "http_method", "method", "value")
	for _, e := range entity {
		qBuilder = qBuilder.Values(e.AppId, e.HttpMethod, e.Method, e.Value)
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
	ctx = sql_metrics.OperationLabelToContext(ctx, "AccessList.UpsertAccessList")

	q := `
	INSERT INTO access_list 
	(app_id, http_method, method, value)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (app_id, http_method, method) 
		DO UPDATE
			SET (value) = (SELECT EXCLUDED.value)
	`
	result, err := r.db.Exec(ctx, q, e.AppId, e.HttpMethod, e.Method, e.Value)
	if err != nil {
		return 0, errors.WithMessagef(err, "exec query %s", q)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, errors.WithMessage(err, "get rows affected")
	}

	return int(rowsAffected), nil
}

func (r AccessList) DeleteAccessList(ctx context.Context, appId int, methods []entity.Method) error {
	ctx = sql_metrics.OperationLabelToContext(ctx, "AccessList.DeleteAccessList")

	methodsClause := squirrel.Or{}
	for _, method := range methods {
		methodsClause = append(methodsClause, squirrel.And{
			squirrel.Eq{"method": method.Method},
			squirrel.Eq{"http_method": method.HttpMethod},
		})
	}

	q, args, err := query.New().
		Delete("access_list").
		Where(squirrel.And{
			squirrel.Eq{
				"app_id": appId,
			},
			methodsClause,
		}).
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
	ctx = sql_metrics.OperationLabelToContext(ctx, "AccessList.DeleteAccessListByAppId")

	q := `
	DELETE FROM access_list
	WHERE app_id = $1
	RETURNING app_id, http_method, method, value
	`
	result := make([]entity.AccessList, 0)
	err := r.db.Select(ctx, &result, q, appId)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}
