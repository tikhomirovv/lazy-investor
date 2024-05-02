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
	"github.com/tikhomirovv/lazy-investor/pkg/config"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
	"os"
)

// Injectors from wire.go:

func InitConfig() (*config.Config, error) {
	string2 := providerApplicationConfigPath()
	configConfig, err := config.NewConfig(string2)
	if err != nil {
		return nil, err
	}
	return configConfig, nil
}

func InitTinkoffService(logger logging.Logger) (*tinkoff.TinkoffService, error) {
	tinkoffConfig := providerTinkoffConfig()
	tinkoffService, err := tinkoff.NewTinkoffService(tinkoffConfig, logger)
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
	configConfig, err := InitConfig()
	if err != nil {
		return nil, err
	}
	zLogger := InitLogger()
	tinkoffService, err := InitTinkoffService(zLogger)
	if err != nil {
		return nil, err
	}
	chartService := chart.NewChartService()
	analyticsService := analytics.NewAnalyticsService()
	applicationApplication := application.NewApplication(configConfig, zLogger, tinkoffService, chartService, analyticsService)
	return applicationApplication, nil
}

// wire.go:

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
