package assembly

import (
	"context"

	"github.com/integration-system/isp-kit/app"
	"github.com/integration-system/isp-kit/bootstrap"
	"github.com/integration-system/isp-kit/cluster"
	"github.com/integration-system/isp-kit/dbrx"
	"github.com/integration-system/isp-kit/dbx"
	"github.com/integration-system/isp-kit/grpc"
	"github.com/integration-system/isp-kit/log"
	rd "github.com/integration-system/isp-lib/v2/redis"
	"github.com/pkg/errors"
	"isp-system-service/conf"
	"isp-system-service/migrations"
	"isp-system-service/redis"
)

type Assembly struct {
	boot         *bootstrap.Bootstrap
	db           *dbrx.Client
	redis        *rd.RxClient
	server       *grpc.Server
	logger       *log.Adapter
	instanceUuid string
}

func New(boot *bootstrap.Bootstrap) (*Assembly, error) {
	dbCli := dbrx.New(dbx.WithMigration(boot.MigrationsDir))
	server := grpc.NewServer()

	localConfig := conf.Local{}
	err := boot.App.Config().Read(&localConfig)
	if err != nil {
		return nil, errors.WithMessage(err, "read local config")
	}

	redisCli := rd.NewRxClient(
		rd.WithInitHandler(func(c *rd.Client, err error) {
			if err != nil {
				boot.App.Logger().Fatal(c.Context(), "redis init", log.Any("err", err))
			}
		}))

	return &Assembly{
		boot:         boot,
		db:           dbCli,
		redis:        redisCli,
		server:       server,
		logger:       boot.App.Logger(),
		instanceUuid: localConfig.InstanceUuid,
	}, nil
}

func (a *Assembly) ReceiveConfig(ctx context.Context, remoteConfig []byte) error {
	var (
		newCfg  conf.Remote
		prevCfg conf.Remote
	)
	err := a.boot.RemoteConfig.Upgrade(remoteConfig, &newCfg, &prevCfg)
	if err != nil {
		a.logger.Fatal(ctx, errors.WithMessage(err, "upgrade remote config"))
	}

	a.redis.ReceiveConfiguration(newCfg.Redis)
	redisCli := redis.NewClient(a.instanceUuid, a.redis)

	migrations.Initialize.SetParams(redisCli, a.boot.MigrationsDir, newCfg.Database.Schema)

	a.logger.SetLevel(newCfg.LogLevel)

	err = a.db.Upgrade(ctx, newCfg.Database)
	if err != nil {
		a.logger.Fatal(ctx, errors.WithMessage(err, "upgrade db client"), log.Any("config", newCfg.Database))
	}

	locator := NewLocator(a.db, redisCli, a.logger)
	handler := locator.Handler(newCfg)
	a.server.Upgrade(handler)

	return nil
}

func (a *Assembly) Runners() []app.Runner {
	eventHandler := cluster.NewEventHandler().
		RemoteConfigReceiver(a)
	return []app.Runner{
		app.RunnerFunc(func(ctx context.Context) error {
			return a.server.ListenAndServe(a.boot.BindingAddress)
		}),
		app.RunnerFunc(func(ctx context.Context) error {
			return a.boot.ClusterCli.Run(ctx, eventHandler)
		}),
	}
}

func (a *Assembly) Closers() []app.Closer {
	return []app.Closer{
		a.boot.ClusterCli,
		app.CloserFunc(func() error {
			a.server.Shutdown()
			return nil
		}),
		a.redis,
		a.db,
	}
}
