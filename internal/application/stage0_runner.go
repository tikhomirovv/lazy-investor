// Package application implements Stage 0 pipeline: fetch candles, compute metrics, build report, send to Telegram.
// Single-flight guard prevents overlapping runs.

package application

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/adapters/report/chart"
	"github.com/tikhomirovv/lazy-investor/internal/application/metrics"
	"github.com/tikhomirovv/lazy-investor/internal/dto"
	"github.com/tikhomirovv/lazy-investor/internal/ports"
)

// runMu guards Stage 0 pipeline so only one run executes at a time.
var runMu sync.Mutex

// RunStage0Once runs the full Stage 0 pipeline once: load candles per instrument, compute metrics, build report, send to Telegram (and optional chart).
// Safe to call from scheduler; concurrent calls are serialized.
func (a *Application) RunStage0Once(ctx context.Context) {
	if !runMu.TryLock() {
		a.logger.Info("Stage 0 run skipped: previous run still in progress")
		return
	}
	defer runMu.Unlock()

	a.logger.Info("Stage 0 pipeline started")
	lookback := time.Duration(a.config.Candles.LookbackDays) * 24 * time.Hour
	to := time.Now()
	from := to.Add(-lookback)

	var rows []InstrumentMetrics
	var firstCandles []dto.Candle
	var firstInstrumentName string

	for _, instConf := range a.config.Instruments {
		instrument, err := a.market.FindInstrument(ctx, instConf.Isin)
		if err != nil {
			a.logger.Warn("FindInstrument failed", "isin", instConf.Isin, "error", err)
			continue
		}
		if instrument == nil {
			a.logger.Warn("instrument not found", "isin", instConf.Isin)
			continue
		}
		candles, err := a.market.GetCandles(ctx, instrument, from, to, ports.Interval1Day)
		if err != nil {
			a.logger.Warn("GetCandles failed", "isin", instConf.Isin, "error", err)
			continue
		}
		if len(candles) == 0 {
			a.logger.Debug("no candles", "isin", instConf.Isin)
			continue
		}
		closes := make([]float64, len(candles))
		volumes := make([]int64, len(candles))
		for i := range candles {
			closes[i] = candles[i].Close
			volumes[i] = candles[i].Volume
		}
		min, max := metrics.MinMax(closes)
		row := InstrumentMetrics{
			Name:        instrument.Name,
			Last:        metrics.Last(closes),
			Change1d:    metrics.PercentChange(closes, 1),
			Change7d:    metrics.PercentChange(closes, 7),
			Change30d:   metrics.PercentChange(closes, 30),
			Min:         min,
			Max:         max,
			AvgVolume:   metrics.AvgVolume(volumes),
			Volatility:  metrics.RealisedVolatility(closes),
			MaxDrawdown: metrics.MaxDrawdown(closes),
		}
		if a.indicators != nil {
			ind := a.indicators.Compute(closes)
			row.SMA20, row.EMA20, row.RSI14 = ind.SMA20, ind.EMA20, ind.RSI14
		}
		rows = append(rows, row)
		if len(firstCandles) == 0 {
			firstCandles = candles
			firstInstrumentName = instrument.Name
		}
	}

	if len(rows) == 0 {
		a.logger.Warn("Stage 0: no data for any instrument, skipping report")
		return
	}

	data := ReportData{AsOf: to, Rows: rows}
	a.logger.Info("Stage 0 report built", "instruments", len(rows))
	// At debug level the full report is visible in logs (FormatForLog = verbose).
	a.logger.Debug("Stage 0 report (full):\n" + FormatForLog(data))

	if a.config.Telegram.Enabled {
		if err := a.telegram.SendMessage(ctx, FormatForTelegram(data)); err != nil {
			a.logger.Error("Telegram SendMessage failed", "error", err)
		}
		if len(firstCandles) > 0 && a.chartSvc != nil {
			var buf bytes.Buffer
			chartInput := candlesToChartInput(firstInstrumentName, firstCandles)
			if err := a.chartSvc.Generate(chartInput, &buf); err != nil {
				a.logger.Warn("chart generate failed", "error", err)
			} else if buf.Len() > 0 {
				if err := a.telegram.SendPhoto(ctx, firstInstrumentName+" (D1)", bytes.NewReader(buf.Bytes())); err != nil {
					a.logger.Error("Telegram SendPhoto failed", "error", err)
				}
			}
		}
	} else {
		a.logger.Debug("Telegram disabled in config, report not sent")
	}
	a.logger.Info("Stage 0 pipeline finished")
}

// candlesToChartInput converts dto candles to chart.Input for PNG generation.
func candlesToChartInput(title string, candles []dto.Candle) *chart.Input {
	times := make([]time.Time, len(candles))
	open, high, low, close := make([]float64, len(candles)), make([]float64, len(candles)), make([]float64, len(candles)), make([]float64, len(candles))
	for i := range candles {
		times[i] = candles[i].Time
		open[i] = candles[i].Open
		high[i] = candles[i].High
		low[i] = candles[i].Low
		close[i] = candles[i].Close
	}
	return &chart.Input{
		Title: title,
		Times: times,
		Open:  open,
		High:  high,
		Low:   low,
		Close: close,
	}
}
