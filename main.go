package main

import (
	"os"

	"isp-system-service/conf"
	"isp-system-service/helper"
	"isp-system-service/model"

	_ "isp-system-service/migrations"

	"github.com/integration-system/isp-lib/backend"
	"github.com/integration-system/isp-lib/database"
	"github.com/integration-system/isp-lib/metric"
	"github.com/integration-system/isp-lib/socket"

	"context"

	"github.com/integration-system/isp-lib/bootstrap"
	"github.com/integration-system/isp-lib/redis"
)

var (
	version = "0.1.0"
	date    = "undefined"
)

func main() {
	bootstrap.
		ServiceBootstrap(&conf.Configuration{}, &conf.RemoteConfig{}).
		OnLocalConfigLoad(onLocalConfigLoad).
		SocketConfiguration(socketConfiguration).
		DeclareMe(routesData).
		OnRemoteConfigReceive(onRemoteConfigReceive).
		OnShutdown(onShutdown).
		Run()
}

/*func ensureRootToken() {
	tokens, err := model.TokenRep.GetTokensByAppId(entity.RootAdminApplicationId)
	if err != nil {
		logger.Error(err)
		return
	}
	if len(tokens) > 0 {
		m, err, _ := controller.GetIdMap(entity.RootAdminApplicationId)
		if err != nil {
			logger.Error(err)
			return
		}
		for _, token := range tokens {
			err = controller.SetIdentityMapForToken(token, m)
			if err != nil {
				logger.Error(err)
				continue
			}
		}
	}
}*/

func socketConfiguration(cfg interface{}) socket.SocketConfiguration {
	appConfig := cfg.(*conf.Configuration)
	return socket.SocketConfiguration{
		Host:   appConfig.ConfigServiceAddress.IP,
		Port:   appConfig.ConfigServiceAddress.Port,
		Secure: false,
		UrlParams: map[string]string{
			"module_name":   appConfig.ModuleName,
			"instance_uuid": appConfig.InstanceUuid,
		},
	}
}

func onShutdown(_ context.Context, _ os.Signal) {
	backend.StopGrpcServer()
	database.Close()
}

func onRemoteConfigReceive(remoteConfig, oldConfig *conf.RemoteConfig) {
	if remoteConfig.RedisAddress.GetAddress() != oldConfig.RedisAddress.GetAddress() {
		rd.InitClient(rd.RedisConfiguration{
			Address:   remoteConfig.RedisAddress,
			DefaultDB: int(rd.ApplicationTokenDb),
		})
	}
	database.InitDb(remoteConfig.DB)
	model.InitDbManager(database.GetDBManager())
	metric.InitCollectors(remoteConfig.Metrics, oldConfig.Metrics)
	metric.InitHttpServer(remoteConfig.Metrics)
	//ensureRootToken()
}

func onLocalConfigLoad(cfg *conf.Configuration) {
	handlers := helper.GetAllHandlers()
	service := backend.GetDefaultService(cfg.ModuleName, handlers...)
	backend.StartBackendGrpcServer(cfg.GrpcInnerAddress, service)
}

func routesData(localConfig interface{}) bootstrap.ModuleInfo {
	cfg := localConfig.(*conf.Configuration)
	return bootstrap.ModuleInfo{
		ModuleName:       cfg.ModuleName,
		ModuleVersion:    version,
		GrpcOuterAddress: cfg.GrpcOuterAddress,
		Handlers:         helper.GetAllHandlers(),
	}
}
