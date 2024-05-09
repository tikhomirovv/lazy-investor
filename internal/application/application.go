package application

import (
	"context"
	"os"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/analytics"
	"github.com/tikhomirovv/lazy-investor/internal/chart"
	"github.com/tikhomirovv/lazy-investor/internal/dto"
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
	instrument, _ := a.tinkoff.GetInstrumentIdByQuery(i.Isin)
	// instrumentId := "BBG004730N88"
	// from := time.Now().Add(-24 * 30 * 12 * time.Hour)
	// to := time.Now().Add(-6 * time.Hour)

	var candles []dto.Candle
	// var dates [][]time.Time = [][] {time.Now().Add(-24 * 30 *12 *time.Hour )}
	dates := [][]time.Time{
		{
			time.Now().Add(-24 * 30 * 24 * time.Hour),
			time.Now().Add(-24 * 30 * 12 * time.Hour),
		},
		{
			time.Now().Add(-24 * 30 * 12 * time.Hour),
			time.Now().Add(-6 * time.Hour),
		}}

	for _, date := range dates {
		ccc, err := a.tinkoff.GetCandles(instrument.Uid, date[0], date[1], tinkoff.CandleIntervalDay)
		if err != nil {
			a.logger.Error("GetCandles error", "error", err)
			return
		}
		a.logger.Debug("GetCandles", "candles", ccc)
		candles = append(candles, ccc...)
	}

	// currentTrend, tc, l, s := a.analytics.AnalyzeTrendByMovingAverage(candles, 30, 80)
	// currentTrend, tc, l := a.analytics.Analyze(candles, 100)
	// a.logger.Debug("Trends", "curr", currentTrend, "trends", tc)
	outFile, err := os.Create(".files/chart" + i.Isin + ".png")
	if err != nil {
		a.logger.Error("Generate chart error", "error", err)
		return
	}
	defer outFile.Close()

	chart := &chart.ChartValues{
		Title:   instrument.Name,
		Candles: candles,
		// Trends: tc,
		EMAs: []dto.EMA{
			analytics.CalculateMovingAverage("EMA 10", candles, 10),
			analytics.CalculateMovingAverage("EMA 50", candles, 50),
			analytics.CalculateMovingAverage("EMA 100", candles, 100),
			analytics.CalculateMovingAverage("EMA 200", candles, 200),
		},
	}
	err = a.chart.Generate(chart, outFile)
	if err != nil {
		a.logger.Error("Generate chart error", "error", err)
	}
}

func (a *Application) Stop() {
	a.tinkoff.Stop()
}
