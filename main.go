package main

import (
	"github.com/integration-system/isp-lib/config/schema"
	"github.com/integration-system/isp-lib/structure"
	log "github.com/integration-system/isp-log"
	"github.com/integration-system/isp-log/stdcodes"
	"isp-system-service/redis"
	"os"

	"isp-system-service/conf"
	"isp-system-service/helper"
	"isp-system-service/model"

	_ "isp-system-service/migrations"

	"context"
	"github.com/integration-system/isp-lib/backend"
	"github.com/integration-system/isp-lib/metric"

	"github.com/integration-system/isp-lib/bootstrap"
)

var (
	version = "0.1.0"
	date    = "undefined"
)

// @title ISP system service
// @version 1.1.2
// @description Сервис управления реестром внешних приложений и токенами аутентификации

// @license.name GNU GPL v3.0

// @host localhost:9003
// @BasePath /api/system
func main() {
	bootstrap.
		ServiceBootstrap(&conf.Configuration{}, &conf.RemoteConfig{}).
		OnLocalConfigLoad(onLocalConfigLoad).
		DefaultRemoteConfigPath(schema.ResolveDefaultConfigPath("default_remote_config.json")).
		SocketConfiguration(socketConfiguration).
		DeclareMe(routesData).
		OnRemoteConfigReceive(onRemoteConfigReceive).
		OnShutdown(onShutdown).
		OnSocketErrorReceive(onRemoteErrorReceive).
		OnConfigErrorReceive(onRemoteConfigErrorReceive).
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

func onRemoteErrorReceive(errorMessage map[string]interface{}) {
	log.WithMetadata(errorMessage).Error(stdcodes.ReceiveErrorFromConfig, "error from config service")
}

func onRemoteConfigErrorReceive(errorMessage string) {
	log.WithMetadata(map[string]interface{}{
		"message": errorMessage,
	}).Error(stdcodes.ReceiveErrorOnGettingConfigFromConfig, "error on getting remote configuration")
}

func socketConfiguration(cfg interface{}) structure.SocketConfiguration {
	appConfig := cfg.(*conf.Configuration)
	return structure.SocketConfiguration{
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
	_ = model.DbClient.Close()
	redis.Client.Close()
}

func onRemoteConfigReceive(remoteConfig, oldConfig *conf.RemoteConfig) {
	redis.Client.ReceiveConfiguration(remoteConfig.Redis)
	model.DbClient.ReceiveConfiguration(remoteConfig.Database)
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
