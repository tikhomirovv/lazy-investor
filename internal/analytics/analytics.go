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

// CalculateMovingAverage calculates the moving average for a slice of candles.
func CalculateMovingAverage(candles []dto.Candle, period int) []float64 {
	ma := make([]float64, len(candles))
	var sum float64
	for i := 0; i < len(candles); i++ {
		sum += candles[i].Close
		if i >= period {
			sum -= candles[i-period].Close
			ma[i] = sum / float64(period)
		} else {
			ma[i] = sum / float64(i+1)
		}
	}
	return ma
}

// AnalyzeTrendByMovingAverage analyzes the trend based on moving average crossovers.
func (a *AnalyticsService) AnalyzeTrendByMovingAverage(candles []dto.Candle, shortPeriod int, longPeriod int) (dto.TrendType, []dto.TrendChange, []float64, []float64) {
	var trendChanges []dto.TrendChange
	if len(candles) < longPeriod {
		return dto.NoTrend, trendChanges, nil, nil // Not enough data to determine the trend
	}
	shortMA := CalculateMovingAverage(candles, shortPeriod)
	longMA := CalculateMovingAverage(candles, longPeriod)
	var currentTrend dto.TrendType = dto.NoTrend

	a.logger.Debug("MA", "short", shortMA, "long", longMA)

	for i := longPeriod; i < len(candles); i++ {
		if shortMA[i] > longMA[i] && shortMA[i-1] <= longMA[i-1] {
			if currentTrend != dto.Uptrend {
				currentTrend = dto.Uptrend
				trendChanges = append(trendChanges, dto.TrendChange{Candle: candles[i], NewTrend: dto.Uptrend})
			}
		} else if shortMA[i] < longMA[i] && shortMA[i-1] >= longMA[i-1] {
			if currentTrend != dto.Downtrend {
				currentTrend = dto.Downtrend
				trendChanges = append(trendChanges, dto.TrendChange{Candle: candles[i], NewTrend: dto.Downtrend})
			}
		}
	}
	return currentTrend, trendChanges, shortMA, longMA
}
