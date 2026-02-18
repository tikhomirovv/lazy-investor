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

// SeriesResult holds a full indicator series aligned to the close series by index.
// Used for features (value/prev/delta, events) and for chart overlays.
type SeriesResult struct {
	Values []float64 // same length as input closes; indices 0..Warmup-1 are zero (not yet ready)
	Ready  bool      // true if at least two valid values exist (for value and prev)
	Warmup int       // minimum candles required (e.g. period for EMA)
}

// IndicatorProvider computes technical indicators. External libs (e.g. cinar/indicator) are used only inside adapters.
type IndicatorProvider interface {
	// Compute returns the last value of each indicator from the close series. Idle/warmup period may yield zeros.
	Compute(closes []float64) IndicatorValues
	// EMA returns the full EMA series for the given period. Values are aligned to closes by index; warmup prefix is zero.
	EMA(closes []float64, period int) SeriesResult
}
