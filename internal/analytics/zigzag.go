package analytics

import (
	"github.com/tikhomirovv/lazy-investor/internal/dto"
)

func CalculateZigZag(candles []dto.Candle, threshold float64) []dto.ZigZagPoint {
	var zigzagPoints []dto.ZigZagPoint
	var lastPivot dto.ZigZagPoint
	var direction string

	for i, candle := range candles {
		if i == 0 {
			lastPivot = dto.ZigZagPoint{Candle: candle, Type: dto.ZigZagHigh, Price: candle.Close}
			zigzagPoints = append(zigzagPoints, lastPivot)
			continue
		}

		priceChangeHigh := (candle.High - lastPivot.Price) / lastPivot.Price
		priceChangeLow := (candle.Low - lastPivot.Price) / lastPivot.Price

		if direction == "" || direction == "down" && priceChangeHigh >= threshold {
			lastPivot = dto.ZigZagPoint{Candle: candle, Type: dto.ZigZagHigh, Price: candle.High}
			zigzagPoints = append(zigzagPoints, lastPivot)
			direction = "up"
		} else if direction == "up" && priceChangeLow <= -threshold {
			lastPivot = dto.ZigZagPoint{Candle: candle, Type: dto.ZigZagLow, Price: candle.Low}
			zigzagPoints = append(zigzagPoints, lastPivot)
			direction = "down"
		}
	}
	return zigzagPoints
}

func ZigZag(candles []dto.Candle, depth int, deviation float64, backstep int) []dto.ZigZagPoint {
	var points []dto.ZigZagPoint
	// Вспомогательная функция для поиска максимального значения
	highest := func(candles []dto.Candle, start, end int) (float64, int) {
		high := candles[start].High
		index := start
		for i := start + 1; i <= end; i++ {
			if candles[i].High > high {
				high = candles[i].High
				index = i
			}
		}
		return high, index
	}
	// Вспомогательная функция для поиска минимального значения
	lowest := func(candles []dto.Candle, start, end int) (float64, int) {
		low := candles[start].Low
		index := start
		for i := start + 1; i <= end; i++ {
			if candles[i].Low < low {
				low = candles[i].Low
				index = i
			}
		}
		return low, index
	}

	// Основной цикл для вычисления точек ZigZag
	for i := depth; i < len(candles); i++ {
		if i+backstep >= len(candles) {
			break
		}
		// Поиск максимума и минимума в пределах глубины
		high, highIndex := highest(candles, i-depth, i)
		low, lowIndex := lowest(candles, i-depth, i)
		// Проверка условий для максимума
		if candles[i].High-high >= deviation && highIndex < i-backstep {
			points = append(points, dto.ZigZagPoint{
				Price:  high,
				Candle: candles[i],
				Type:   dto.ZigZagHigh,
			})
		}
		// Проверка условий для минимума
		if low-candles[i].Low >= deviation && lowIndex < i-backstep {
			points = append(points, dto.ZigZagPoint{
				Price:  low,
				Candle: candles[i],
				Type:   dto.ZigZagLow,
			})
		}
	}
	return points
}
