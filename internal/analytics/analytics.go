package analytics

import "github.com/tikhomirovv/lazy-investor/internal/dto"

type AnalyticsService struct {
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{}
}

func (a *AnalyticsService) Analyze(candles []dto.Candle) []dto.TrendChange {
	var trendChanges []dto.TrendChange
	var currentTrend dto.TrendType = dto.NoTrend
	var lastHigh float64 = 0
	var lastLow float64 = 0

	for _, candle := range candles {
		if currentTrend == dto.NoTrend {
			trendChanges = append(trendChanges, dto.TrendChange{Candle: candle, NewTrend: dto.NoTrend})
			currentTrend = dto.Uptrend
			lastHigh = candle.High
			lastLow = candle.Low
		} else if currentTrend == dto.Uptrend {
			if candle.High > lastHigh {
				lastHigh = candle.High
			}
			if candle.Close < lastLow {
				currentTrend = dto.Downtrend
				trendChanges = append(trendChanges, dto.TrendChange{Candle: candle, NewTrend: dto.Downtrend})
				lastLow = candle.Low
			}
		} else if currentTrend == dto.Downtrend {
			if candle.Low < lastLow {
				lastLow = candle.Low
			}
			if candle.Close > lastHigh {
				currentTrend = dto.Uptrend
				trendChanges = append(trendChanges, dto.TrendChange{Candle: candle, NewTrend: dto.Uptrend})
				lastHigh = candle.High
			}
		}
	}
	return trendChanges
}
