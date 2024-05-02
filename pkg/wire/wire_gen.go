// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/tikhomirovv/lazy-investor/internal/analytics"
	"github.com/tikhomirovv/lazy-investor/internal/application"
	"github.com/tikhomirovv/lazy-investor/internal/chart"
	"github.com/tikhomirovv/lazy-investor/internal/tinkoff"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
	"os"
)

// Injectors from wire.go:

func InitTinkoffService(logger logging.Logger) (*tinkoff.TinkoffService, error) {
	config := providerTinkoffConfig()
	tinkoffService, err := tinkoff.NewTinkoffService(config, logger)
	if err != nil {
		return nil, err
	}
	return tinkoffService, nil
}

func InitLogger() *logging.ZLogger {
	zLogger := logging.NewLogger()
	return zLogger
}

func InitApplication() (*application.Application, error) {
	zLogger := InitLogger()
	tinkoffService, err := InitTinkoffService(zLogger)
	if err != nil {
		return nil, err
	}
	chartService := chart.NewChartService()
	analyticsService := analytics.NewAnalyticsService()
	applicationApplication := application.NewApplication(zLogger, tinkoffService, chartService, analyticsService)
	return applicationApplication, nil
}

// wire.go:

func providerTinkoffConfig() tinkoff.Config {
	return tinkoff.Config{
		AppName: os.Getenv("APP_NAME"),
		Host:    os.Getenv("TINKOFF_API_HOST"),
		Token:   os.Getenv("TINKOFF_API_TOKEN"),
	}

}
