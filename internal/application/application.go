// Package application provides the main app orchestration. Run/Stop and wiring of adapters.
package application

import (
	"context"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/adapters/report/chart"
	"github.com/tikhomirovv/lazy-investor/internal/ports"
	"github.com/tikhomirovv/lazy-investor/pkg/config"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

// Application holds config and adapters; Run executes the Stage 0 pipeline (scheduler → data → report → telegram), Stop cleans up.
type Application struct {
	config    *config.Config
	logger    logging.Logger
	market    ports.MarketDataProvider
	chartSvc  *chart.Service
	telegram  ports.TelegramNotifier
}

// NewApplication constructs the app from config and adapters (used by Wire).
func NewApplication(
	cfg *config.Config,
	logger logging.Logger,
	market ports.MarketDataProvider,
	chartSvc *chart.Service,
	telegram ports.TelegramNotifier,
) *Application {
	return &Application{
		config:   cfg,
		logger:   logger,
		market:   market,
		chartSvc: chartSvc,
		telegram: telegram,
	}
}

// Run runs the Stage 0 loop: optional run-on-start, then periodic runs by scheduler interval until ctx is done.
func (a *Application) Run(ctx context.Context) {
	a.logger.Info("Application started (Stage 0 pipeline)")
	if a.config.Scheduler.RunOnStart {
		a.RunStage0Once(ctx)
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
