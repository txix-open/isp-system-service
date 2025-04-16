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

type Token struct {
	db db.DB
}

func NewToken(db db.DB) Token {
	return Token{
		db: db,
	}
}

func (r Token) SaveToken(ctx context.Context, token string, appId int, expireTime int) (*entity.Token, error) {
	q := `
	INSERT INTO token
	(token, app_id, expire_time)
	VALUES ($1, $2, $3)
	RETURNING token, app_id, expire_time, created_at
`
	result := entity.Token{}
	err := r.db.SelectRow(ctx, &result, q, token, appId, expireTime)
	if err != nil {
		return nil, errors.WithMessage(err, "select row db")
	}

	return &result, nil
}

func (r Token) GetTokenById(ctx context.Context, token string) (*entity.Token, error) {
	q := `
	SELECT token, app_id, expire_time, created_at
	FROM token
	WHERE token = $1
`
	result := entity.Token{}
	err := r.db.SelectRow(ctx, &result, q, token)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil //nolint:nilnil
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}

func (r Token) GetTokenByAppIdList(ctx context.Context, appIdList []int) ([]entity.Token, error) {
	q, args, err := query.New().
		Select("token", "app_id", "expire_time", "created_at").
		From("token").
		Where(squirrel.Eq{"app_id": appIdList}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, errors.WithMessage(err, "build query")
	}

	result := make([]entity.Token, 0)
	err = r.db.Select(ctx, &result, q, args...)
	if err != nil {
		return nil, errors.WithMessagef(err, "exec query %s", q)
	}

	return result, nil
}

func (r Token) DeleteToken(ctx context.Context, tokens []string) (int, error) {
	q, args, err := query.New().
		Delete("token").
		Where(squirrel.Eq{"token": tokens}).
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

func (r Token) AuthDataByToken(ctx context.Context, token string) (*entity.AuthData, error) {
	q := `
SELECT system_id, domain_id, application_group_id, app_id, token.expire_time, token.created_at
FROM token
         LEFT JOIN application
                   ON token.app_id = application.id
         LEFT JOIN application_group
                   ON application.application_group_id = application_group.id
         LEFT JOIN domain
                   ON application_group.domain_id = domain.id
WHERE token = $1
`
	result := entity.AuthData{}
	err := r.db.SelectRow(ctx, &result, q, token)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrTokenNotFound
	case err != nil:
		return nil, errors.WithMessagef(err, "exec query %s", q)
	default:
		return &result, nil
	}
}
