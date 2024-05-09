package analytics

import (
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/dto"
)

// CalculateMovingAverage calculates the moving average for a slice of candles.
func CalculateMovingAverage(name string, candles []dto.Candle, period int) dto.EMA {
	dates := make([]time.Time, len(candles))
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
		dates[i] = candles[i].Time
	}
	return dto.EMA{
		Name:   name,
		Dates:  dates,
		Values: ma,
	}
}
