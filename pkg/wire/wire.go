//go:build wireinject
// +build wireinject

// Package wire provides dependency injection via Google Wire. Run: make wire or wire gen .
package wire

import (
	"os"

	"github.com/google/wire"
	"github.com/tikhomirovv/lazy-investor/internal/adapters/indicators/indicator"
	"github.com/tikhomirovv/lazy-investor/internal/adapters/marketdata/tinkoff"
	"github.com/tikhomirovv/lazy-investor/internal/adapters/report/chart"
	"github.com/tikhomirovv/lazy-investor/internal/adapters/telegram"
	"github.com/tikhomirovv/lazy-investor/internal/application"
	"github.com/tikhomirovv/lazy-investor/internal/application/exportcandles"
	"github.com/tikhomirovv/lazy-investor/internal/ports"
	"github.com/tikhomirovv/lazy-investor/pkg/config"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

const configPath = "./config.yml"

func providerConfigPath() string {
	return configPath
}

func providerTinkoffConfig() tinkoff.Config {
	return tinkoff.Config{
		AppName: os.Getenv("APP_NAME"),
		Host:    os.Getenv("TINKOFF_API_HOST"),
		Token:   os.Getenv("TINKOFF_API_TOKEN"),
	}
}

// providerTelegramConfig returns Telegram adapter config from env. Empty token/chatID = no-op mode.
func providerTelegramConfig() telegram.Config {
	return telegram.Config{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID: os.Getenv("TELEGRAM_CHAT_ID"),
	}
}

// InitConfig returns config from YAML file.
func InitConfig() (*config.Config, error) {
	panic(wire.Build(
		providerConfigPath,
		config.NewConfig,
	))
}

// InitLogger returns the default zerolog-based logger.
func InitLogger() *logging.ZLogger {
	panic(wire.Build(
		logging.NewLogger,
	))
}

// InitTinkoffService builds Tinkoff market data adapter.
func InitTinkoffService(logger logging.Logger) (*tinkoff.Service, error) {
	panic(wire.Build(
		providerTinkoffConfig,
		tinkoff.NewService,
	))
}

// InitApplication builds the main application with all adapters (market data, chart, telegram, indicators, export candles).
func InitApplication() (*application.Application, error) {
	panic(wire.Build(
		providerConfigPath,
		config.NewConfig,
		providerTinkoffConfig,
		providerTelegramConfig,
		wire.Bind(new(logging.Logger), new(*logging.ZLogger)),
		wire.Bind(new(ports.MarketDataProvider), new(*tinkoff.Service)),
		wire.Bind(new(ports.TelegramNotifier), new(*telegram.Service)),
		wire.Bind(new(ports.IndicatorProvider), new(*indicator.Service)),
		logging.NewLogger,
		tinkoff.NewService,
		telegram.NewService,
		chart.NewService,
		indicator.NewService,
		exportcandles.NewService,
		application.NewApplication,
	))
}
