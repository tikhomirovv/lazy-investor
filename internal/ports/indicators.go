// Package ports defines interfaces (contracts) between app/domain and adapters.
// IndicatorProvider: compute technical indicators from OHLC/close series. Implemented by adapters (e.g. cinar/indicator).

package ports

// IndicatorValues holds the last value of each indicator (0 = not available).
// All indicators are computed from the same close series; caller passes []float64.
type IndicatorValues struct {
	SMA20 float64 // Simple Moving Average, period 20
	EMA20 float64 // Exponential Moving Average, period 20
	RSI14 float64 // Relative Strength Index, period 14 (0–100)
}

// IndicatorProvider computes technical indicators. External libs (e.g. cinar/indicator) are used only inside adapters.
type IndicatorProvider interface {
	// Compute returns the last value of each indicator from the close series. Idle/warmup period may yield zeros.
	Compute(closes []float64) IndicatorValues
}
