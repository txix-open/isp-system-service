package transaction

import (
	"context"
	"isp-system-service/service/baseline"

	"isp-system-service/repository"
	"isp-system-service/service"

	"github.com/txix-open/isp-kit/db"
)

type Manager struct {
	db db.Transactional
}

func NewManager(db db.Transactional) *Manager {
	return &Manager{
		db: db,
	}
}

type accessListSetOneTx struct {
	repository.AccessList
}

func (m Manager) AccessListSetOneTx(ctx context.Context, msgTx func(ctx context.Context, tx service.AccessListSetOneTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		accessListRepository := repository.NewAccessList(tx)
		return msgTx(ctx, accessListSetOneTx{
			AccessList: accessListRepository,
		})
	})
}

type accessListSetListTx struct {
	repository.AccessList
}

func (m Manager) AccessListSetListTx(ctx context.Context, msgTx func(ctx context.Context, tx service.AccessListSetListTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		accessListRep := repository.NewAccessList(tx)
		return msgTx(ctx, accessListSetListTx{
			AccessList: accessListRep,
		})
	})
}

type applicationDeleteTx struct {
	repository.Application
}

func (m Manager) ApplicationDeleteTx(ctx context.Context, msgTx func(ctx context.Context, tx service.ApplicationDeleteTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		applicationRep := repository.NewApplication(tx)
		return msgTx(ctx, applicationDeleteTx{
			Application: applicationRep,
		})
	})
}

type tokenCreateTx struct {
	repository.Token
}

func (m Manager) TokenCreateTx(ctx context.Context, msgTx func(ctx context.Context, tx service.TokenCreateTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		tokenRep := repository.NewToken(tx)
		return msgTx(ctx, tokenCreateTx{
			Token: tokenRep,
		})
	})
}

type tokenRevokeTx struct {
	repository.Token
}

func (m Manager) TokenRevokeTx(ctx context.Context, msgTx func(ctx context.Context, tx service.TokenRevokeTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		tokenRep := repository.NewToken(tx)
		return msgTx(ctx, tokenRevokeTx{
			Token: tokenRep,
		})
	})
}

type baselineTx struct {
	repository.Locker
	repository.Domain
	repository.AppGroup
	repository.Application
	repository.AccessList
	repository.Token
}

func (m Manager) BaselineTx(ctx context.Context, txTx func(ctx context.Context, tx baseline.Transaction) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		return txTx(ctx, baselineTx{
			Locker:      repository.NewLocker(tx),
			Domain:      repository.NewDomain(tx),
			AppGroup:    repository.NewAppGroup(tx),
			Application: repository.NewApplication(tx),
			AccessList:  repository.NewAccessList(tx),
			Token:       repository.NewToken(tx),
		})
	})
}
