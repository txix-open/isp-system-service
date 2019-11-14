package model

import (
	"github.com/integration-system/isp-lib/database"
	log "github.com/integration-system/isp-log"
	"isp-system-service/log_code"
)

var (
	DbClient = database.NewRxDbClient(
		database.WithSchemaEnsuring(),
		database.WithSchemaAutoInjecting(),
		database.WithMigrationsEnsuring(),
		database.WithInitializingErrorHandler(func(err *database.ErrorEvent) {
			log.Error(log_code.ErrorDatabaseClient, err)
		}))

	AppRep        = AppRepository{rxClient: DbClient}
	DomainRep     = DomainRepository{rxClient: DbClient}
	ServiceRep    = ServiceRepository{rxClient: DbClient}
	SystemRep     = SystemRepository{rxClient: DbClient}
	TokenRep      = TokenRepository{rxClient: DbClient}
	AccessListRep = AccessListRepository{rxClient: DbClient}
)
