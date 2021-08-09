package main

import (
	"context"
	"os"

	"github.com/integration-system/isp-lib/v2/backend"
	"github.com/integration-system/isp-lib/v2/bootstrap"
	"github.com/integration-system/isp-lib/v2/config/schema"
	"github.com/integration-system/isp-lib/v2/metric"
	"github.com/integration-system/isp-lib/v2/structure"
	log "github.com/integration-system/isp-log"
	"github.com/integration-system/isp-log/stdcodes"
	"isp-system-service/conf"
	_ "isp-system-service/docs"
	"isp-system-service/helper"
	"isp-system-service/migrations"
	"isp-system-service/model"
	"isp-system-service/redis"
)

var version = "0.1.0"

// @title ISP system service
// @version 1.1.2
// @description Сервис управления реестром внешних приложений и токенами аутентификации

// @license.name GNU GPL v3.0

// @host localhost:9000
// @BasePath /api/system

//go:generate swag init --parseDependency
//go:generate rm -f docs/swagger.json
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
		SubscribeBroadcastEvent(bootstrap.ListenRestartEvent()).
		Run()
}

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
			"module_name": appConfig.ModuleName,
		},
	}
}

func onShutdown(_ context.Context, _ os.Signal) {
	backend.StopGrpcServer()
	_ = model.DbClient.Close()
	redis.Client.Close()
}

func onRemoteConfigReceive(remoteConfig, oldConfig *conf.RemoteConfig) {
	migrations.DatabaseConfig = remoteConfig.Database
	redis.Client.ReceiveConfiguration(remoteConfig.Redis)
	model.DbClient.ReceiveConfiguration(remoteConfig.Database)
	metric.InitCollectors(remoteConfig.Metrics, oldConfig.Metrics)
	metric.InitHttpServer(remoteConfig.Metrics)
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
