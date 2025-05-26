package assembly

import (
	"context"

	"github.com/txix-open/isp-kit/rc"

	"isp-system-service/conf"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/app"
	"github.com/txix-open/isp-kit/bootstrap"
	"github.com/txix-open/isp-kit/cluster"
	"github.com/txix-open/isp-kit/dbrx"
	"github.com/txix-open/isp-kit/dbx"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/log"
)

type Assembly struct {
	boot   *bootstrap.Bootstrap
	db     *dbrx.Client
	server *grpc.Server
	logger *log.Adapter
}

func New(boot *bootstrap.Bootstrap) (*Assembly, error) {
	logger := boot.App.Logger()
	dbCli := dbrx.New(logger, dbx.WithMigrationRunner(boot.MigrationsDir, logger))
	server := grpc.NewServer()
	return &Assembly{
		boot:   boot,
		db:     dbCli,
		server: server,
		logger: logger,
	}, nil
}

func (a *Assembly) ReceiveConfig(ctx context.Context, remoteConfig []byte) error {
	newCfg, _, err := rc.Upgrade[conf.Remote](a.boot.RemoteConfig, remoteConfig)
	if err != nil {
		a.boot.Fatal(errors.WithMessage(err, "upgrade remote config"))
	}

	a.logger.SetLevel(newCfg.LogLevel)

	err = a.db.Upgrade(ctx, newCfg.Database)
	if err != nil {
		a.boot.Fatal(errors.WithMessage(err, "upgrade db client"))
	}

	locator := NewLocator(a.db, a.logger)
	config := locator.Config(newCfg)

	err = config.Baseline.Do(a.boot.App.Context())
	if err != nil {
		a.boot.Fatal(errors.WithMessage(err, "run baseline"))
	}

	a.server.Upgrade(config.Handler)

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
		a.db,
	}
}
