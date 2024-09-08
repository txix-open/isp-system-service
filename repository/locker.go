package repository

import (
	"context"
	"hash/fnv"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/db"
	"github.com/txix-open/isp-kit/metrics/sql_metrics"
)

const (
	prefix = "msp-service-template"
)

type Locker struct {
	db db.DB
}

func NewLocker(db db.DB) Locker {
	return Locker{db: db}
}

func (l Locker) TryLock(ctx context.Context, key string) (bool, error) {
	ctx = sql_metrics.OperationLabelToContext(ctx, "Locker.TryLock")

	hash := fnv.New32a()
	_, err := hash.Write([]byte(prefix + key))
	if err != nil {
		return false, errors.WithMessage(err, "generate hash")
	}
	sum := hash.Sum32()

	acquired := false
	err = l.db.SelectRow(ctx, &acquired, "SELECT pg_try_advisory_xact_lock($1)", sum)
	if err != nil {
		return false, errors.WithMessage(err, "pg try acquire advisory lock")
	}
	return acquired, nil
}
