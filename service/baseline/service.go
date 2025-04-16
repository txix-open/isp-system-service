package baseline

import (
	"context"
	"isp-system-service/conf"
	"isp-system-service/entity"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/log"
)

const adminAppId = 1

type Transaction interface {
	CreateDomain(ctx context.Context, name string, desc string, systemId int) (*entity.Domain, error)
	CreateAppGroup(ctx context.Context, name string, desc string, domainId int) (*entity.AppGroup, error)
	CreateApplication(ctx context.Context, id int, name string, desc string, appGroupId int, appType string) (*entity.Application, error)
	UpsertAccessList(ctx context.Context, e entity.AccessList) (int, error)
	SaveToken(ctx context.Context, token string, appId int, expireTime int) (*entity.Token, error)
	GetTokenById(ctx context.Context, token string) (*entity.Token, error)
	TryLock(ctx context.Context, key string) (bool, error)
}

type TxRunner interface {
	BaselineTx(ctx context.Context, tx func(ctx context.Context, tx Transaction) error) error
}

type Service struct {
	cfg      conf.Baseline
	txRunner TxRunner
	logger   log.Logger
}

func NewService(cfg conf.Baseline, txRunner TxRunner, logger log.Logger) Service {
	return Service{
		cfg:      cfg,
		txRunner: txRunner,
		logger:   logger,
	}
}

func (s Service) Do(ctx context.Context) error {
	ctx = log.ToContext(ctx, log.String("worker", "baseline"))
	if s.cfg.InitialAdminUiToken == "" {
		s.logger.Info(ctx, "initial admin ui token is empty, skip baseline")
		return nil
	}

	err := s.txRunner.BaselineTx(ctx, s.transaction)
	if err != nil {
		return errors.WithMessage(err, "run baseline transaction")
	}

	return nil
}

func (s Service) transaction(ctx context.Context, tx Transaction) error {
	locked, err := tx.TryLock(ctx, "isp-system-service.baseline")
	if err != nil {
		return errors.WithMessage(err, "try lock isp-system-service.baseline")
	}
	if !locked {
		s.logger.Info(ctx, "baseline is locked, skip baseline")
		return nil
	}

	token, err := tx.GetTokenById(ctx, s.cfg.InitialAdminUiToken)
	if err != nil {
		return errors.WithMessage(err, "get token by id")
	}
	if token != nil {
		s.logger.Info(ctx, "initial admin ui token is exist, skip baseline")
		return nil
	}

	s.logger.Info(ctx, "initial admin ui token is empty, run baseline")

	domain, err := tx.CreateDomain(ctx, "root", "", 1)
	if err != nil {
		return errors.WithMessage(err, "create domain")
	}

	service, err := tx.CreateAppGroup(ctx, "rootService", "", domain.Id)
	if err != nil {
		return errors.WithMessage(err, "create service")
	}

	app, err := tx.CreateApplication(ctx, adminAppId, "admin", "", service.Id, "SYSTEM")
	if err != nil {
		return errors.WithMessage(err, "create application")
	}

	_, err = tx.SaveToken(ctx, s.cfg.InitialAdminUiToken, app.Id, -1)
	if err != nil {
		return errors.WithMessage(err, "save token")
	}

	_, err = tx.UpsertAccessList(ctx, entity.AccessList{
		AppId:  app.Id,
		Method: "admin/auth/login",
		Value:  true,
	})
	if err != nil {
		return errors.WithMessage(err, "upsert access list")
	}

	s.logger.Info(ctx, "baseline done")

	return nil
}
