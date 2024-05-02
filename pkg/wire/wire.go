//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/tikhomirovv/lazy-investor/internal/application"
	"github.com/tikhomirovv/lazy-investor/internal/tinkoff"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

func InitTinkoffService(logger logging.Logger) (*tinkoff.TinkoffService, error) {
	// appName := os.Getenv("APP_NAME")
	// host := os.Getenv("TINKOFF_API_HOST")
	// token := os.Getenv("TINKOFF_API_TOKEN")
	wire.Build(
		tinkoff.NewTinkoffService,
	)
	return &tinkoff.TinkoffService{}, nil
}

func InitApplication() (*application.Application, error) {
	wire.Build(
		logging.NewLogger,
		wire.Bind(new(logging.Logger), new(*logging.ZLogger)),
		InitTinkoffService,
		application.NewApplication,
	)
	return &application.Application{}, nil
}
