// Package features: unit tests for EMA feature computation (determinism, warmup, events).
package features

import (
	"testing"

	"github.com/tikhomirovv/lazy-investor/internal/ports"
)

// fakeIndicatorProvider returns configurable EMA series for testing (no external lib).
type fakeIndicatorProvider struct {
	emaFunc func(closes []float64, period int) ports.SeriesResult
}

func (f *fakeIndicatorProvider) Compute(closes []float64) ports.IndicatorValues {
	return ports.IndicatorValues{}
}

func (f *fakeIndicatorProvider) EMA(closes []float64, period int) ports.SeriesResult {
	if f.emaFunc != nil {
		return f.emaFunc(closes, period)
	}
	return ports.SeriesResult{Values: make([]float64, len(closes)), Warmup: period}
}

// TestComputeEMAFeatures_Determinism verifies that same input produces same output.
func TestComputeEMAFeatures_Determinism(t *testing.T) {
	closes := []float64{100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120}
	provider := &fakeIndicatorProvider{
		emaFunc: func(c []float64, period int) ports.SeriesResult {
			sr := ports.SeriesResult{Warmup: period, Values: make([]float64, len(c)), Ready: len(c) >= period && len(c) >= 2}
			for i := period - 1; i < len(c); i++ {
				sr.Values[i] = float64(100 + i)
			}
			return sr
		},
	}
	a := ComputeEMAFeatures(provider, closes)
	b := ComputeEMAFeatures(provider, closes)
	if a.EMA20.Value != b.EMA20.Value || a.EMA20.Ready != b.EMA20.Ready {
		t.Errorf("determinism failed: two calls gave different EMA20")
	}
	if a.EMA100.Value != b.EMA100.Value {
		t.Errorf("determinism failed: two calls gave different EMA100")
	}
}

// TestComputeEMAFeatures_Warmup verifies ready=false when not enough data.
func TestComputeEMAFeatures_Warmup(t *testing.T) {
	// provider returns Ready only when len >= period and at least 2 values
	provider := &fakeIndicatorProvider{
		emaFunc: func(closes []float64, period int) ports.SeriesResult {
			sr := ports.SeriesResult{Warmup: period, Values: make([]float64, len(closes))}
			if len(closes) < period {
				return sr
			}
			for i := period - 1; i < len(closes); i++ {
				sr.Values[i] = float64(i)
			}
			sr.Ready = len(closes) >= period && (len(closes)-period+1) >= 2
			return sr
		},
	}
	// period-1 candles: not ready
	short := make([]float64, 19)
	for i := range short {
		short[i] = 100 + float64(i)
	}
	out := ComputeEMAFeatures(provider, short)
	if out.EMA20.Ready {
		t.Errorf("EMA20 should not be ready for 19 candles (period 20)")
	}
	// period candles and 2+ values: ready
	long := make([]float64, 25)
	for i := range long {
		long[i] = 100 + float64(i)
	}
	out2 := ComputeEMAFeatures(provider, long)
	if !out2.EMA20.Ready {
		t.Errorf("EMA20 should be ready for 25 candles")
	}
}

// TestComputeEMAFeatures_PriceCrossedUp verifies price_crossed_up event when prevClose <= prev and lastClose > value.
func TestComputeEMAFeatures_PriceCrossedUp(t *testing.T) {
	// Last two closes: 98, 102. EMA last two: 99, 100. So prevClose 98 <= 99, lastClose 102 > 100 -> cross up
	closes := make([]float64, 21)
	for i := range closes {
		closes[i] = 90 + float64(i)
	}
	closes[19] = 98
	closes[20] = 102
	provider := &fakeIndicatorProvider{
		emaFunc: func(c []float64, period int) ports.SeriesResult {
			sr := ports.SeriesResult{Warmup: period, Values: make([]float64, len(c)), Ready: true}
			sr.Values[19] = 99
			sr.Values[20] = 100
			return sr
		},
	}
	out := ComputeEMAFeatures(provider, closes)
	found := false
	for _, e := range out.EMA20.Events {
		if e == "price_crossed_up_ema20" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected price_crossed_up_ema20 in events, got %v", out.EMA20.Events)
	}
}

// TestComputeEMAFeatures_PriceCrossedDown verifies price_crossed_down when prevClose >= prev and lastClose < value.
func TestComputeEMAFeatures_PriceCrossedDown(t *testing.T) {
	closes := make([]float64, 21)
	for i := range closes {
		closes[i] = 100 + float64(i)
	}
	closes[19] = 102
	closes[20] = 98
	provider := &fakeIndicatorProvider{
		emaFunc: func(c []float64, period int) ports.SeriesResult {
			sr := ports.SeriesResult{Warmup: period, Values: make([]float64, len(c)), Ready: true}
			sr.Values[19] = 99
			sr.Values[20] = 100
			return sr
		},
	}
	out := ComputeEMAFeatures(provider, closes)
	found := false
	for _, e := range out.EMA20.Events {
		if e == "price_crossed_down_ema20" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected price_crossed_down_ema20 in events, got %v", out.EMA20.Events)
	}
}

// TestComputeEMAFeatures_NilProvider returns zero struct.
func TestComputeEMAFeatures_NilProvider(t *testing.T) {
	out := ComputeEMAFeatures(nil, []float64{1, 2, 3})
	if out.EMA20.Ready || out.EMA100.Ready {
		t.Errorf("nil provider should yield not ready")
	}
}

// TestComputeEMAFeatures_ShortCloses returns zero when len(closes) < 2.
func TestComputeEMAFeatures_ShortCloses(t *testing.T) {
	provider := &fakeIndicatorProvider{}
	out := ComputeEMAFeatures(provider, []float64{100})
	if out.EMA20.Ready {
		t.Errorf("single close should yield not ready")
	}
}
