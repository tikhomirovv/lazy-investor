// Package application provides the main app orchestration. Run/Stop and wiring of adapters.
package application

import (
	"context"

	"github.com/tikhomirovv/lazy-investor/internal/adapters/marketdata/tinkoff"
	"github.com/tikhomirovv/lazy-investor/internal/adapters/report/chart"
	"github.com/tikhomirovv/lazy-investor/pkg/config"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

// Application holds config and adapters; Run executes the pipeline (none yet), Stop cleans up.
type Application struct {
	config  *config.Config
	logger  logging.Logger
	tinkoff *tinkoff.Service
	chart   *chart.Service
}

// NewApplication constructs the app from config and adapters (used by Wire).
func NewApplication(
	cfg *config.Config,
	logger logging.Logger,
	tinkoffSvc *tinkoff.Service,
	chartSvc *chart.Service,
) *Application {
	return &Application{
		config:  cfg,
		logger:  logger,
		tinkoff: tinkoffSvc,
		chart:   chartSvc,
	}
}

// Run runs the main loop (scheduler → data → features → …). No pipeline yet; placeholder.
func (a *Application) Run(ctx context.Context) {
	a.logger.Info("Application started (no pipeline yet)")
	<-ctx.Done()
}

// Stop closes external connections (e.g. Tinkoff client).
func (a *Application) Stop() {
	a.tinkoff.Stop()
}
