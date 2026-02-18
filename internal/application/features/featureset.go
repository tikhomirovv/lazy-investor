// Package features provides deterministic feature computation from indicator series.
// Used by the report and (later) by snapshot/LLM. All logic is pure: same input → same output.
package features

// EMAFeature holds the computed feature set for one EMA(period): value, prev, delta, ready, and events.
// See ComputeEMAFeatures and the doc-comment in ema.go for the full specification.
type EMAFeature struct {
	Value  float64   // EMA on the last closed candle
	Prev   float64   // EMA on the previous closed candle
	Delta  float64   // Value - Prev
	Ready  bool      // true if at least two valid EMA points exist
	Events []string  // e.g. price_above_ema20, price_crossed_up_ema20, price_crossed_down_ema20
}

// EMAFeatureSet holds features for EMA20, EMA100, and the combined trend filter ema20_above_ema100.
type EMAFeatureSet struct {
	EMA20             EMAFeature // features for period 20
	EMA100            EMAFeature // features for period 100
	EMA20AboveEMA100  bool       // true when EMA20 > EMA100 (trend filter)
	EMA20AboveReady   bool       // true when both EMA20 and EMA100 are ready
	CombinedEvents    []string   // ema20_crossed_above_ema100, ema20_crossed_below_ema100
}
