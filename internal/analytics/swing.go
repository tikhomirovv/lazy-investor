package analytics

import (
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/dto"
)

// Функция для поиска swing highs и swing lows
func FindSwings(candles []dto.Candle, n int) ([]time.Time, []time.Time) {
	var swingHighs []time.Time
	var swingLows []time.Time

	for i := n; i < len(candles)-n; i++ {
		isSwingHigh := true
		isSwingLow := true
		for j := -n; j <= n; j++ {
			if j != 0 {
				if candles[i].High <= candles[i+j].High {
					isSwingHigh = false
				}
				if candles[i].Low >= candles[i+j].Low {
					isSwingLow = false
				}
			}
		}
		if isSwingHigh {
			swingHighs = append(swingHighs, candles[i].Time)
		}
		if isSwingLow {
			swingLows = append(swingLows, candles[i].Time)
		}
	}
	return swingHighs, swingLows
}
