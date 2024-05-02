//go:build wireinject
// +build wireinject

package wire

import (
	"os"

	"github.com/google/wire"
	"github.com/tikhomirovv/lazy-investor/internal/application"
	"github.com/tikhomirovv/lazy-investor/internal/tinkoff"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

func providerTinkoffConfig() tinkoff.Config {
	return tinkoff.Config{
		AppName: os.Getenv("APP_NAME"),
		Host:    os.Getenv("TINKOFF_API_HOST"),
		Token:   os.Getenv("TINKOFF_API_TOKEN"),
	}

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
		InitLogger,
		wire.Bind(new(logging.Logger), new(*logging.ZLogger)),
		InitTinkoffService,
		application.NewApplication,
	)
	return &application.Application{}, nil
}
