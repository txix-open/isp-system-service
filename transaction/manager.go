package transaction

import (
	"context"

	"github.com/integration-system/isp-kit/db"
	"isp-system-service/repository"
	"isp-system-service/service"
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

func (m Manager) AccessListSetOneTx(ctx context.Context, msgTx func(ctx context.Context, tx service.IAccessListSetOneTx) error) error {
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

func (m Manager) AccessListSetListTx(ctx context.Context, msgTx func(ctx context.Context, tx service.IAccessListSetListTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		accessListRep := repository.NewAccessList(tx)
		return msgTx(ctx, accessListSetListTx{
			AccessList: accessListRep,
		})
	})
}

type applicationDeleteTx struct {
	repository.AccessList
	repository.Application
	repository.Token
}

func (m Manager) ApplicationDeleteTx(ctx context.Context, msgTx func(ctx context.Context, tx service.IApplicationDeleteTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		accessListRep := repository.NewAccessList(tx)
		applicationRep := repository.NewApplication(tx)
		tokenRep := repository.NewToken(tx)
		return msgTx(ctx, applicationDeleteTx{
			AccessList:  accessListRep,
			Application: applicationRep,
			Token:       tokenRep,
		})
	})
}

type tokenCreateTx struct {
	repository.Token
}

func (m Manager) TokenCreateTx(ctx context.Context, msgTx func(ctx context.Context, tx service.ITokenCreateTx) error) error {
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

func (m Manager) TokenRevokeTx(ctx context.Context, msgTx func(ctx context.Context, tx service.ITokenRevokeTx) error) error {
	return m.db.RunInTransaction(ctx, func(ctx context.Context, tx *db.Tx) error {
		tokenRep := repository.NewToken(tx)
		return msgTx(ctx, tokenRevokeTx{
			Token: tokenRep,
		})
	})
}
