//go:build wireinject
// +build wireinject

package wire

import (
	"os"

	"github.com/google/wire"
	"github.com/tikhomirovv/lazy-investor/internal/application"
	"github.com/tikhomirovv/lazy-investor/internal/services"
	"github.com/tikhomirovv/lazy-investor/pkg/config"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

const (
	ConfigPath = "./config.yml"
)

func providerApplicationConfigPath() string {
	return ConfigPath
}

func providerTinkoffConfig() services.TinkoffConfig {
	return services.TinkoffConfig{
		AppName: os.Getenv("APP_NAME"),
		Host:    os.Getenv("TINKOFF_API_HOST"),
		Token:   os.Getenv("TINKOFF_API_TOKEN"),
	}

}

func InitConfig() (*config.Config, error) {
	panic(wire.Build(
		providerApplicationConfigPath,
		config.NewConfig,
	))
}

func InitTinkoffService(logger logging.Logger) (*services.TinkoffService, error) {
	wire.Build(
		providerTinkoffConfig,
		services.NewTinkoffService,
	)
	return &services.TinkoffService{}, nil
}

func InitLogger() *logging.ZLogger {
	wire.Build(
		logging.NewLogger,
	)
	return &logging.ZLogger{}
}

func InitApplication() (*application.Application, error) {
	wire.Build(
		InitConfig,
		InitLogger,
		wire.Bind(new(logging.Logger), new(*logging.ZLogger)),
		InitTinkoffService,
		services.NewChartService,
		services.NewStrategyService,
		application.NewApplication,
	)
	return &application.Application{}, nil
}
