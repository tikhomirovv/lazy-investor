//go:build wireinject
// +build wireinject

package wire

import (
	"os"

	"github.com/google/wire"
	"github.com/tikhomirovv/lazy-investor/internal/application"
	"github.com/tikhomirovv/lazy-investor/internal/chart"
	"github.com/tikhomirovv/lazy-investor/internal/tinkoff"
	"github.com/tikhomirovv/lazy-investor/pkg/config"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

const (
	ConfigPath = "./config.yml"
)

func providerApplicationConfigPath() string {
	return ConfigPath
}

func providerTinkoffConfig() tinkoff.Config {
	return tinkoff.Config{
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

func InitTinkoffService(logger logging.Logger) (*tinkoff.TinkoffService, error) {
	wire.Build(
		providerTinkoffConfig,
		tinkoff.NewTinkoffService,
	)
	return &tinkoff.TinkoffService{}, nil
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
		chart.NewChartService,
		application.NewApplication,
	)
	return &application.Application{}, nil
}
