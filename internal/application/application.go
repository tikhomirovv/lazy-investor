package application

import (
	"context"
	"fmt"
	"os"

	"github.com/tikhomirovv/lazy-investor/internal/analytics"
	"github.com/tikhomirovv/lazy-investor/internal/dto"
	"github.com/tikhomirovv/lazy-investor/internal/services"
	"github.com/tikhomirovv/lazy-investor/pkg/config"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

type Application struct {
	config   *config.Config
	logger   logging.Logger
	tinkoff  *services.TinkoffService
	chart    *services.ChartService
	strategy *services.StrategyService
}

func NewApplication(
	config *config.Config,
	logger logging.Logger,
	tinkoff *services.TinkoffService,
	chart *services.ChartService,
	strategy *services.StrategyService,
) *Application {
	return &Application{
		config:   config,
		logger:   logger,
		tinkoff:  tinkoff,
		chart:    chart,
		strategy: strategy,
	}
}

func (a *Application) Run(ctx context.Context) {

	var instruments []*dto.Instrument
	for _, i := range a.config.Instruments {
		instrument, _ := a.tinkoff.GetInstrumentIdByQuery(i.Isin)
		instruments = append(instruments, instrument)
		// go a.analyse(instrument)
	}

	a.strategy.Test(instruments)
}

func (a *Application) analyse(instrument *dto.Instrument) error {
	candles, err := a.tinkoff.GetCandles(instrument)
	if err != nil {
		return fmt.Errorf("application.analyse: %w", err)
	}
	// currentTrend, tc, l, s := a.analytics.AnalyzeTrendByMovingAverage(candles, 30, 80)
	// currentTrend, tc, l := a.analytics.Analyze(candles, 100)
	// a.logger.Debug("Trends", "curr", currentTrend, "trends", tc)
	outFile, err := os.Create(".files/chart" + string(instrument.Isin) + ".png")
	if err != nil {
		a.logger.Error("Create chart file", "error", err)
		return fmt.Errorf("application.analyse: Create chart file: %w", err)
	}
	defer outFile.Close()

	swings := analytics.FindSwings(candles, 2)
	currentTrend, trendChanges := analytics.GetTrends(swings)
	// zz := analytics.CalculateZigZag(candles, 0.02)
	// zz := analytics.ZigZag(candles, 14, 0.01, 2)
	chart := &services.ChartValues{
		Title:   instrument.Name,
		Candles: candles,
		Trends:  trendChanges,
		// EMAs: []dto.EMA{
		// 	analytics.CalculateMovingAverage("EMA 10", candles, 10),
		// 	analytics.CalculateMovingAverage("EMA 50", candles, 50),
		// 	analytics.CalculateMovingAverage("EMA 100", candles, 100),
		// 	analytics.CalculateMovingAverage("EMA 200", candles, 200),
		// },
		Swings: swings,
		// ZigZags: zz,
	}
	a.logger.Info("Current trend", "trend", currentTrend.String(), "tc", trendChanges)
	// a.logger.Info("ZZ", "zz", zz)
	err = a.chart.Generate(chart, outFile)
	if err != nil {
		a.logger.Error("Generate chart error", "error", err)
		return fmt.Errorf("application.analyse: Generate chart error: %w", err)
	}

	return nil
}

func (a *Application) Stop() {
	a.tinkoff.Stop()
}
