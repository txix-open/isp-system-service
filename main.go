package main

import (
	"github.com/txix-open/isp-kit/bootstrap"
	"github.com/txix-open/isp-kit/shutdown"
	"isp-system-service/assembly"
	"isp-system-service/conf"
	"isp-system-service/routes"
)

var (
	version = "1.0.0"
)

//	@title			isp-system-service
//	@version		1.0.0
//	@description	Сервис управления реестром внешних приложений и токенами аутентификации

//	@license.name	GNU GPL v3.0

//	@host		localhost:9000
//	@BasePath	/api/system

//go:generate swag init --parseDependency
//go:generate rm -f docs/swagger.json docs/docs.go
func main() {
	boot := bootstrap.New(version, conf.Remote{}, routes.EndpointDescriptors())
	app := boot.App
	logger := app.Logger()

	assembly, err := assembly.New(boot)
	if err != nil {
		logger.Fatal(app.Context(), err)
	}
	app.AddRunners(assembly.Runners()...)
	app.AddClosers(assembly.Closers()...)

	shutdown.On(func() {
		logger.Info(app.Context(), "starting shutdown")
		app.Shutdown()
		logger.Info(app.Context(), "shutdown completed")
	})

	err = app.Run()
	if err != nil {
		boot.Fatal(err)
	}
}
