package pkg

import "math"

// Функция для вычисления среднего значения массива
func Average(values []float64) float64 {
	var sum float64
	for _, value := range values {
		sum += value
	}
	// Возвращаем среднее значение
	return sum / float64(len(values))
}

// Функция для вычисления стандартного отклонения (риска) для массива доходностей
func StandardDeviation(values []float64) float64 {
	// Вычисляем среднее значение доходности
	mean := Average(values)
	var variance float64
	// Для каждого значения в массиве вычисляем отклонение от среднего, возводим в квадрат и суммируем
	for _, value := range values {
		variance += math.Pow(value-mean, 2)
	}
	// Стандартное отклонение — это корень из дисперсии
	return math.Sqrt(variance / float64(len(values)))
}

// Функция для вычисления коэффициента Шарпа
func SharpeRatio(returns []float64, riskFreeRate float64) float64 {
	// Расчет средней доходности портфеля
	meanReturn := Average(returns)
	// Расчет стандартного отклонения (риска)
	risk := StandardDeviation(returns)
	// Коэффициент Шарпа: (доходность - безрисковая ставка) / риск
	return (meanReturn - riskFreeRate) / risk
}
