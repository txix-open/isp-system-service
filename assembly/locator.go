package assembly

import (
	"github.com/txix-open/isp-kit/db"
	"github.com/txix-open/isp-kit/grpc/endpoint"
	"github.com/txix-open/isp-kit/grpc/isp"
	"github.com/txix-open/isp-kit/log"
	"isp-system-service/conf"
	"isp-system-service/controller"
	"isp-system-service/repository"
	"isp-system-service/routes"
	"isp-system-service/service"
	"isp-system-service/service/secure"
	"isp-system-service/transaction"
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

func (l Locator) Handler(cfg conf.Remote) isp.BackendServiceServer {
	txManager := transaction.NewManager(l.db)
	accessListRep := repository.NewAccessList(l.db)
	applicationRep := repository.NewApplication(l.db)
	domainRep := repository.NewDomain(l.db)
	serviceRep := repository.NewService(l.db)
	tokenRep := repository.NewToken(l.db)

	secureService := secure.NewService(tokenRep, accessListRep)
	accessListService := service.NewAccessList(txManager, accessListRep, applicationRep)
	applicationService := service.NewApplication(txManager, applicationRep, domainRep, serviceRep, tokenRep)
	domainService := service.NewDomain(domainRep)
	serviceService := service.NewService(domainRep, serviceRep)

	jwtService := service.NewTokenSource()
	tokenService := service.NewToken(jwtService, applicationService, txManager,
		applicationRep, domainRep, serviceRep, tokenRep,
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
	mapper := endpoint.DefaultWrapper(l.logger, endpoint.BodyLogger(l.logger))
	server := routes.Handler(mapper, c)
	return server
}
