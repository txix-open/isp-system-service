package assembly

import (
	"isp-system-service/conf"
	"isp-system-service/controller"
	"isp-system-service/repository"
	"isp-system-service/routes"
	"isp-system-service/service"
	"isp-system-service/service/baseline"
	"isp-system-service/service/secure"
	"isp-system-service/transaction"

	"github.com/txix-open/isp-kit/db"
	"github.com/txix-open/isp-kit/grpc"
	"github.com/txix-open/isp-kit/grpc/endpoint"
	"github.com/txix-open/isp-kit/grpc/endpoint/grpclog"
	"github.com/txix-open/isp-kit/log"
)

type DB interface {
	db.DB
	db.Transactional
}

type Locator struct {
	db     DB
	logger log.Logger
}

func NewLocator(db DB, logger log.Logger) Locator {
	return Locator{
		db:     db,
		logger: logger,
	}
}

type Config struct {
	Handler  *grpc.Mux
	Baseline baseline.Service
}

func (l Locator) Config(cfg conf.Remote) Config {
	txManager := transaction.NewManager(l.db)
	accessListRep := repository.NewAccessList(l.db)
	applicationRep := repository.NewApplication(l.db)
	domainRep := repository.NewDomain(l.db)
	appGroupRep := repository.NewAppGroup(l.db)
	tokenRep := repository.NewToken(l.db)

	secureService := secure.NewService(tokenRep, accessListRep)
	accessListService := service.NewAccessList(txManager, accessListRep, applicationRep)
	applicationService := service.NewApplication(txManager, applicationRep, domainRep, appGroupRep, tokenRep)
	domainService := service.NewDomain(domainRep)
	serviceService := service.NewService(domainRep, appGroupRep)

	jwtService := service.NewTokenSource()
	tokenService := service.NewToken(jwtService, applicationService, txManager,
		applicationRep, domainRep, appGroupRep, tokenRep,
	)

	secureController := controller.NewSecure(secureService)
	accessListController := controller.NewAccessList(accessListService)
	applicationController := controller.NewApplication(applicationService)
	domainController := controller.NewDomain(domainService)
	serviceController := controller.NewService(serviceService)
	tokenController := controller.NewToken(tokenService)

	c := routes.Controllers{
		Secure:      secureController,
		AccessList:  accessListController,
		Domain:      domainController,
		Service:     serviceController,
		Application: applicationController,
		Token:       tokenController,
	}
	mapper := endpoint.DefaultWrapper(l.logger, grpclog.Log(l.logger, true))
	server := routes.Handler(mapper, c)

	baselineService := baseline.NewService(cfg.Baseline, txManager, l.logger)
	return Config{
		Handler:  server,
		Baseline: baselineService,
	}
}
