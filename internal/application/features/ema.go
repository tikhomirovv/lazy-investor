// Package features: EMA feature computation (see featureset.go for types).
package features

import (
	"github.com/tikhomirovv/lazy-investor/internal/ports"
)

// ComputeEMAFeatures builds the full EMA feature set for closes using the given indicator provider.
// It is deterministic: same closes and provider → same EMAFeatureSet.
//
// --- Feature specification: EMA(period) ---
//
// Purpose:
//   - EMA measures smoothed price and acts as a trend filter. We use it for value/prev/delta
//     and for events: price above/below EMA, and price crossing up/down through EMA.
//
// Inputs:
//   - close series from fully closed candles only (e.g. daily close). Same length as candle count.
//
// Parameters:
//   - period: smoothing period (we use 20 and 100 in Stage 1).
//
// Warmup / Lookback:
//   - ready = true only when len(closes) >= period and the provider returns at least two valid
//     EMA values (so we have value and prev). Otherwise ready = false.
//
// Output fields (per period):
//   - value: EMA at the last closed candle.
//   - prev: EMA at the previous closed candle.
//   - delta: value - prev.
//   - ready: as above.
//   - events: see Events semantics below.
//
// Events semantics (deterministic, no LLM):
//   - price_above_ema{period}: last close > value (price is above EMA).
//   - price_crossed_up_ema{period}: previous close <= prev EMA and last close > value (cross up).
//   - price_crossed_down_ema{period}: previous close >= prev EMA and last close < value (cross down).
//
// Combined feature (ema20_above_ema100):
//   - value (bool): EMA20 > EMA100.
//   - ready: ema20.ready && ema100.ready.
//   - events: ema20_crossed_above_ema100, ema20_crossed_below_ema100 when the two EMAs cross.
//
// Edge cases:
//   - If ready = false: value, prev, delta are zero; events slice is empty.
//   - If len(closes) < 2: always ready = false.
//   - Provider returns empty or short series: we set ready = false and zero values.
//
// Tests:
//   - Determinism: same closes → same EMAFeatureSet.
//   - Warmup: period-1 candles → ready false; period candles and at least 2 values → ready true.
//   - Cross events: synthetic series where close crosses EMA from below/above produce correct events.

// ComputeEMAFeatures builds the full EMA feature set (EMA20, EMA100, ema20_above_ema100) from closes.
// Uses provider.EMA(closes, period) for each period. Deterministic.
func ComputeEMAFeatures(provider ports.IndicatorProvider, closes []float64) EMAFeatureSet {
	out := EMAFeatureSet{}
	if provider == nil || len(closes) < 2 {
		return out
	}
	out.EMA20 = buildEMAFeature(provider, closes, 20, "ema20")
	out.EMA100 = buildEMAFeature(provider, closes, 100, "ema100")
	// Combined: ema20 above ema100 and cross events
	if out.EMA20.Ready && out.EMA100.Ready {
		out.EMA20AboveReady = true
		out.EMA20AboveEMA100 = out.EMA20.Value > out.EMA100.Value
		prevAbove := out.EMA20.Prev > out.EMA100.Prev
		if !prevAbove && out.EMA20AboveEMA100 {
			out.CombinedEvents = append(out.CombinedEvents, "ema20_crossed_above_ema100")
		}
		if prevAbove && !out.EMA20AboveEMA100 {
			out.CombinedEvents = append(out.CombinedEvents, "ema20_crossed_below_ema100")
		}
	}
	return out
}

func buildEMAFeature(provider ports.IndicatorProvider, closes []float64, period int, namePrefix string) EMAFeature {
	f := EMAFeature{}
	sr := provider.EMA(closes, period)
	if !sr.Ready || len(sr.Values) < 2 {
		return f
	}
	n := len(sr.Values)
	lastClose := closes[n-1]
	prevClose := closes[n-2]
	f.Value = sr.Values[n-1]
	f.Prev = sr.Values[n-2]
	f.Delta = f.Value - f.Prev
	f.Ready = true
	if lastClose > f.Value {
		f.Events = append(f.Events, "price_above_"+namePrefix)
	}
	if prevClose <= f.Prev && lastClose > f.Value {
		f.Events = append(f.Events, "price_crossed_up_"+namePrefix)
	}
	if prevClose >= f.Prev && lastClose < f.Value {
		f.Events = append(f.Events, "price_crossed_down_"+namePrefix)
	}
	return f
}
