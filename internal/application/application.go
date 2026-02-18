// Package application provides the main app orchestration. Run/Stop and wiring of adapters.
package application

import (
	"context"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/adapters/report/chart"
	"github.com/tikhomirovv/lazy-investor/internal/application/exportcandles"
	"github.com/tikhomirovv/lazy-investor/internal/ports"
	"github.com/tikhomirovv/lazy-investor/pkg/config"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

// Application holds config and adapters; Run executes the Stage 0 pipeline (scheduler → data → report → telegram), Stop cleans up.
type Application struct {
	config        *config.Config
	logger        logging.Logger
	market        ports.MarketDataProvider
	chartSvc      *chart.Service
	telegram      ports.TelegramNotifier
	indicators    ports.IndicatorProvider
	exportCandles *exportcandles.Service
}

// NewApplication constructs the app from config and adapters (used by Wire). indicators may be nil.
func NewApplication(
	cfg *config.Config,
	logger logging.Logger,
	market ports.MarketDataProvider,
	chartSvc *chart.Service,
	telegram ports.TelegramNotifier,
	indicators ports.IndicatorProvider,
	exportCandles *exportcandles.Service,
) *Application {
	return &Application{
		config:        cfg,
		logger:        logger,
		market:        market,
		chartSvc:      chartSvc,
		telegram:      telegram,
		indicators:    indicators,
		exportCandles: exportCandles,
	}
}

// Run runs the Stage 0 loop: optional run-on-start, then periodic runs by scheduler interval until ctx is done.
// If telegram.handleCommands is true, also starts a goroutine that listens for /candles etc.
func (a *Application) Run(ctx context.Context) {
	a.logger.Info("Application started (Stage 0 pipeline)")
	if a.config.Scheduler.RunOnStart {
		a.RunStage0Once(ctx)
	}
	if a.config.Telegram.HandleCommands && a.exportCandles != nil {
		go a.runTelegramCommandListener(ctx)
	}
	interval := time.Duration(a.config.Scheduler.IntervalSeconds) * time.Second
	if interval <= 0 {
		interval = time.Hour
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.RunStage0Once(ctx)
		}
	}
}

// Stop closes external connections (e.g. market data provider).
func (a *Application) Stop() {
	a.market.Stop()
}
