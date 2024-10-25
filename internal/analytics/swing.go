package analytics

import (
	"github.com/tikhomirovv/lazy-investor/internal/dto"
)

// TODO: экстремум 2-го порядка
// https://youtu.be/dCZ6R25_t8A?si=dCS29G5Ljhc0XRX7&t=1011
// как это сделано? (это 3 и 4 порядки? то есть сначала 2 порядки, а потом результат по той же логике прогнать!)
// https://youtu.be/dCZ6R25_t8A?si=esdBEISQmHo-to5K&t=5203

// Функция для поиска swing highs и swing lows
// `n“ определяет количество свечей с каждой стороны от текущей свечи,
// которые будут использоваться для определения,
// является ли текущая свеча swing high или swing low
func FindSwings(candles []dto.Candle, n int) []dto.Swing {
	var swings []dto.Swing
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
		if isSwingHigh || isSwingLow {
			swing := dto.Swing{
				Candle: candles[i],
				Period: n,
			}
			if isSwingHigh {
				swing.Type = dto.SwingHigh
			}
			if isSwingLow {
				swing.Type = dto.SwingLow
			}
			swings = append(swings, swing)
		}
	}
	return swings
}
