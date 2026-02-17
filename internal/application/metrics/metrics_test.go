package metrics

import (
	"math"
	"testing"
)

func TestLast(t *testing.T) {
	if got := Last(nil); got != 0 {
		t.Errorf("Last(nil) = %v, want 0", got)
	}
	if got := Last([]float64{}); got != 0 {
		t.Errorf("Last([]) = %v, want 0", got)
	}
	if got := Last([]float64{1, 2, 3}); got != 3 {
		t.Errorf("Last([1,2,3]) = %v, want 3", got)
	}
}

func TestPercentChange(t *testing.T) {
	// 100 -> 110: +10%
	closes := []float64{100, 110}
	if got := PercentChange(closes, 1); math.Abs(got-10) > 0.01 {
		t.Errorf("PercentChange([100,110], 1) = %v, want 10", got)
	}
	// not enough data
	if got := PercentChange([]float64{100}, 1); got != 0 {
		t.Errorf("PercentChange([100], 1) = %v, want 0", got)
	}
	if got := PercentChange(closes, 0); got != 0 {
		t.Errorf("PercentChange(..., 0) = %v, want 0", got)
	}
	// 50 -> 25: -50%
	closes2 := []float64{50, 25}
	if got := PercentChange(closes2, 1); math.Abs(got-(-50)) > 0.01 {
		t.Errorf("PercentChange([50,25], 1) = %v, want -50", got)
	}
}

func TestMinMax(t *testing.T) {
	min, max := MinMax(nil)
	if min != 0 || max != 0 {
		t.Errorf("MinMax(nil) = %v, %v; want 0, 0", min, max)
	}
	min, max = MinMax([]float64{3, 1, 2})
	if min != 1 || max != 3 {
		t.Errorf("MinMax([3,1,2]) = %v, %v; want 1, 3", min, max)
	}
}

func TestAvgVolume(t *testing.T) {
	if got := AvgVolume(nil); got != 0 {
		t.Errorf("AvgVolume(nil) = %v, want 0", got)
	}
	if got := AvgVolume([]int64{10, 20, 30}); got != 20 {
		t.Errorf("AvgVolume([10,20,30]) = %v, want 20", got)
	}
}

func TestDailyReturns(t *testing.T) {
	closes := []float64{100, 110, 99}
	ret := DailyReturns(closes)
	if len(ret) != 2 {
		t.Fatalf("len(DailyReturns) = %d, want 2", len(ret))
	}
	// (110-100)/100 = 0.1, (99-110)/110 ≈ -0.1
	if math.Abs(ret[0]-0.1) > 1e-9 {
		t.Errorf("DailyReturns[0] = %v, want 0.1", ret[0])
	}
	if math.Abs(ret[1]-(-11.0/110.0)) > 1e-9 {
		t.Errorf("DailyReturns[1] = %v", ret[1])
	}
}

func TestStdDev(t *testing.T) {
	// 2, 4, 4, 4, 5, 5, 7, 9: mean 5, variance = (9+1+1+1+0+0+4+16)/8 = 4, std = 2
	vals := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	if got := StdDev(vals); math.Abs(got-2) > 0.01 {
		t.Errorf("StdDev(...) = %v, want 2", got)
	}
	if StdDev(nil) != 0 || StdDev([]float64{1}) != 0 {
		t.Error("StdDev with len < 2 should be 0")
	}
}

func TestRealisedVolatility(t *testing.T) {
	// constant closes -> zero returns -> zero vol
	if got := RealisedVolatility([]float64{1, 1, 1}); got != 0 {
		t.Errorf("RealisedVolatility(const) = %v, want 0", got)
	}
	// 100, 101, 100: returns +0.01, -0.0099...; non-zero std
	closes := []float64{100, 101, 100}
	vol := RealisedVolatility(closes)
	if vol <= 0 || vol > 0.1 {
		t.Errorf("RealisedVolatility([100,101,100]) = %v, expect small positive", vol)
	}
}

func TestMaxDrawdown(t *testing.T) {
	// peak 100, then 80: drawdown 20%
	closes := []float64{100, 80}
	if got := MaxDrawdown(closes); math.Abs(got-0.2) > 0.001 {
		t.Errorf("MaxDrawdown([100,80]) = %v, want 0.2", got)
	}
	// 100, 120, 90: peak 120, trough 90 -> 30/120 = 0.25
	closes2 := []float64{100, 120, 90}
	if got := MaxDrawdown(closes2); math.Abs(got-0.25) > 0.001 {
		t.Errorf("MaxDrawdown([100,120,90]) = %v, want 0.25", got)
	}
	if MaxDrawdown(nil) != 0 || MaxDrawdown([]float64{1}) != 0 {
		t.Error("MaxDrawdown with len < 2 should be 0")
	}
}
