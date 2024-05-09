package analytics

import (
	"github.com/tikhomirovv/lazy-investor/internal/dto"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

type AnalyticsService struct {
	logger logging.Logger
}

func NewAnalyticsService(logger logging.Logger) *AnalyticsService {
	return &AnalyticsService{
		logger: logger,
	}
}

type Direction string

const (
	DirectionUp   Direction = "up"
	DirectionDown Direction = "down"
)

// AnalyzeTrendByMovingAverage analyzes the trend based on moving average crossovers.
func (a *AnalyticsService) AnalyzeTrendByMovingAverage(candles []dto.Candle, shortPeriod int, longPeriod int) (dto.TrendType, []dto.TrendChange, []float64, []float64) {
	var trendChanges []dto.TrendChange
	if len(candles) < longPeriod {
		return dto.NoTrend, trendChanges, nil, nil // Not enough data to determine the trend
	}
	shortMA := CalculateMovingAverage("shortMA", candles, shortPeriod)
	longMA := CalculateMovingAverage("longMA", candles, longPeriod)
	var currentTrend dto.TrendType = dto.NoTrend

	a.logger.Debug("MA", "short", shortMA, "long", longMA)

	for i := longPeriod; i < len(candles); i++ {
		if shortMA.Values[i] > longMA.Values[i] && shortMA.Values[i-1] <= longMA.Values[i-1] {
			if currentTrend != dto.Uptrend {
				currentTrend = dto.Uptrend
				trendChanges = append(trendChanges, dto.TrendChange{Candle: candles[i], NewTrend: dto.Uptrend})
			}
		} else if shortMA.Values[i] < longMA.Values[i] && shortMA.Values[i-1] >= longMA.Values[i-1] {
			if currentTrend != dto.Downtrend {
				currentTrend = dto.Downtrend
				trendChanges = append(trendChanges, dto.TrendChange{Candle: candles[i], NewTrend: dto.Downtrend})
			}
		}
	}
	return currentTrend, trendChanges, shortMA.Values, longMA.Values
}

func (a *AnalyticsService) Analyze(candles []dto.Candle, period int) (dto.TrendType, []dto.TrendChange, []float64) {
	var trendChanges []dto.TrendChange
	var currentTrend dto.TrendType = dto.NoTrend
	var directionUp bool = true
	shortMA := CalculateMovingAverage("", candles, period)
	for i, candle := range candles {
		if i == 0 {
			continue
		}
		prev := shortMA.Values[i-1]
		curr := shortMA.Values[i]
		if curr > prev && !directionUp { // разворот вверх
			directionUp = true
			currentTrend = dto.Uptrend
			trendChanges = append(trendChanges, dto.TrendChange{Candle: candle, NewTrend: currentTrend})
		}
		if curr < prev && directionUp { // разворот вниз
			directionUp = false
			currentTrend = dto.Downtrend
			trendChanges = append(trendChanges, dto.TrendChange{Candle: candle, NewTrend: currentTrend})
		}
	}
	return currentTrend, trendChanges, shortMA.Values
}
