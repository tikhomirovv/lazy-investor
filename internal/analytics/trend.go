package analytics

import "github.com/tikhomirovv/lazy-investor/internal/dto"

func GetTrends(swings []dto.Swing) (dto.TrendType, []dto.TrendChange) {
	var trendChanges []dto.TrendChange
	var currentTrend dto.TrendType = dto.TrendNo
	var prevHigh, lastHigh, prevLow, lastLow float64
	addTrendChange := func(trend dto.TrendType, swing dto.Swing) {
		trendChanges = append(trendChanges, dto.TrendChange{
			Swing: swing,
			Trend: trend,
		})
	}
	for _, swing := range swings {
		if swing.Type == dto.SwingHigh {
			// UpTrend
			// if swing.Candle.High > lastHigh && prevLow > lastLow {
			// 	// change?
			// 	if currentTrend != dto.TrendUp {
			// 		currentTrend = dto.TrendUp
			// 		addTrendChange(currentTrend, swing)
			// 	}
			// } else if currentTrend == dto.TrendUp {
			// 	// no trend
			// 	// if currentTrend != dto.TrendNo {
			// 	currentTrend = dto.TrendNo
			// 	addTrendChange(currentTrend, swing)
			// 	// }
			// }
			prevHigh = lastHigh
			lastHigh = swing.Candle.High
		} else {
			// DownTrend
			// if swing.Candle.Low < lastLow && prevHigh < lastHigh {
			// 	// change?
			// 	if currentTrend != dto.TrendDown {
			// 		currentTrend = dto.TrendDown
			// 		addTrendChange(currentTrend, swing)
			// 	}
			// } else if currentTrend == dto.TrendDown {
			// 	// no trend
			// 	// if currentTrend != dto.TrendNo {
			// 	currentTrend = dto.TrendNo
			// 	addTrendChange(currentTrend, swing)
			// 	// }
			// }
			prevLow = lastLow
			lastLow = swing.Candle.Low
		}

		// uptrend
		if lastHigh > prevHigh && lastLow > prevLow {
			if currentTrend != dto.TrendUp {
				currentTrend = dto.TrendUp
				addTrendChange(currentTrend, swing)
			}
		} else if lastLow < prevLow && lastHigh < prevHigh {
			if currentTrend != dto.TrendDown {
				currentTrend = dto.TrendDown
				addTrendChange(currentTrend, swing)
			}
		} else if currentTrend != dto.TrendNo {
			currentTrend = dto.TrendNo
			addTrendChange(currentTrend, swing)
		}
	}
	return currentTrend, trendChanges
}
