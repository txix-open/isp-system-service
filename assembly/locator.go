package assembly

import (
	"github.com/integration-system/isp-kit/db"
	"github.com/integration-system/isp-kit/grpc/endpoint"
	"github.com/integration-system/isp-kit/grpc/isp"
	"github.com/integration-system/isp-kit/log"
	"isp-system-service/conf"
	"isp-system-service/controller"
	"isp-system-service/redis"
	"isp-system-service/repository"
	"isp-system-service/routes"
	"isp-system-service/service"
	"isp-system-service/transaction"
)

type DB interface {
	db.DB
	db.Transactional
}

type Locator struct {
	db     DB
	redis  redis.Client
	logger log.Logger
}

func NewLocator(db DB, redisCli redis.Client, logger log.Logger) Locator {
	return Locator{
		db:     db,
		redis:  redisCli,
		logger: logger,
	}
}

func (l Locator) Handler(cfg conf.Remote) isp.BackendServiceServer {
	txManager := transaction.NewManager(l.db)
	accessListRep := repository.NewAccessList(l.db)
	applicationRep := repository.NewApplication(l.db)
	domainRep := repository.NewDomain(l.db)
	serviceRep := repository.NewService(l.db)
	tokenRep := repository.NewToken(l.db)

	accessListService := service.NewAccessList(l.redis, txManager, accessListRep, applicationRep)
	applicationService := service.NewApplication(l.redis, txManager, applicationRep, domainRep, serviceRep, tokenRep)
	domainService := service.NewDomain(domainRep)
	serviceService := service.NewService(domainRep, serviceRep)

	jwtService := service.NewJwt(cfg.ApplicationSecret)
	tokenService := service.NewToken(l.redis, cfg.DefaultTokenExpireTime, jwtService, applicationService, txManager,
		applicationRep, domainRep, serviceRep, tokenRep,
	)

	accessListController := controller.NewAccessList(accessListService)
	applicationController := controller.NewApplication(applicationService)
	domainController := controller.NewDomain(domainService)
	serviceController := controller.NewService(serviceService)
	tokenController := controller.NewToken(tokenService)

	c := routes.Controllers{
		AccessList:  accessListController,
		Domain:      domainController,
		Service:     serviceController,
		Application: applicationController,
		Token:       tokenController,
	}
	mapper := endpoint.DefaultWrapper(l.logger, endpoint.BodyLogger(l.logger))
	server := routes.Handler(mapper, c)
	return server
}
