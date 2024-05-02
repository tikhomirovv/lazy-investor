package application

import (
	"context"
	"os"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/analytics"
	"github.com/tikhomirovv/lazy-investor/internal/chart"
	"github.com/tikhomirovv/lazy-investor/internal/tinkoff"
	"github.com/tikhomirovv/lazy-investor/pkg/config"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

type Application struct {
	config    *config.Config
	logger    logging.Logger
	tinkoff   *tinkoff.TinkoffService
	chart     *chart.ChartService
	analytics *analytics.AnalyticsService
}

func NewApplication(config *config.Config, logger logging.Logger, tinkoff *tinkoff.TinkoffService, chart *chart.ChartService, analytics *analytics.AnalyticsService) *Application {
	return &Application{
		config:    config,
		logger:    logger,
		tinkoff:   tinkoff,
		chart:     chart,
		analytics: analytics,
	}
}

func (a *Application) Run(ctx context.Context) {
	for _, i := range a.config.Instruments {
		go a.analyse(i)
	}
}
func (a *Application) analyse(i config.InstConf) {
	uid, _ := a.tinkoff.GetInstrumentIdByQuery(i.Isin)
	// instrumentId := "BBG004730N88"
	instrumentId := uid
	from := time.Now().Add(-24 * 30 * 5 * time.Hour)
	to := time.Now().Add(-6 * time.Hour)

	candles, err := a.tinkoff.GetCandles(instrumentId, from, to, tinkoff.CandleIntervalDay)
	if err != nil {
		a.logger.Error("GetCandles error", "error", err)
		return
	}
	a.logger.Debug("Candles", "i", instrumentId, "candles", candles)

	tc := a.analytics.Analyze(candles)
	a.logger.Debug("Trends", "trends", tc)

	outFile, err := os.Create(".files/chart" + i.Isin + ".png")
	if err != nil {
		a.logger.Error("Generate chart error", "error", err)
		return
	}
	defer outFile.Close()
	err = a.chart.Generate(candles, tc, outFile)
	if err != nil {
		a.logger.Error("Generate chart error", "error", err)
	}
}

func (a *Application) Stop() {
	a.tinkoff.Stop()
}
