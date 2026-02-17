// Package metrics provides pure, deterministic functions for Stage 0 report metrics.
// No external dependencies; operates on slices (closes, volumes). Used by report builder.
package metrics

import "math"

// Last returns the last close price. Returns 0 if empty.
func Last(closes []float64) float64 {
	if len(closes) == 0 {
		return 0
	}
	return closes[len(closes)-1]
}

// PercentChange returns (last - closes[len-periodsAgo]) / closes[len-periodsAgo] * 100.
// periodsAgo 1 = change vs previous bar. Returns 0 if not enough data or divisor 0.
func PercentChange(closes []float64, periodsAgo int) float64 {
	if periodsAgo <= 0 || len(closes) < periodsAgo+1 {
		return 0
	}
	prev := closes[len(closes)-1-periodsAgo]
	if prev == 0 {
		return 0
	}
	return (closes[len(closes)-1] - prev) / prev * 100
}

// MinMax returns min and max of closes. Returns (0,0) if empty.
func MinMax(closes []float64) (min, max float64) {
	if len(closes) == 0 {
		return 0, 0
	}
	min, max = closes[0], closes[0]
	for _, c := range closes[1:] {
		if c < min {
			min = c
		}
		if c > max {
			max = c
		}
	}
	return min, max
}

// AvgVolume returns the mean of volumes. Returns 0 if empty.
func AvgVolume(volumes []int64) float64 {
	if len(volumes) == 0 {
		return 0
	}
	var sum int64
	for _, v := range volumes {
		sum += v
	}
	return float64(sum) / float64(len(volumes))
}

// DailyReturns computes (close[i] - close[i-1]) / close[i-1] for i=1..len-1. One fewer element than closes.
func DailyReturns(closes []float64) []float64 {
	if len(closes) < 2 {
		return nil
	}
	out := make([]float64, 0, len(closes)-1)
	for i := 1; i < len(closes); i++ {
		if closes[i-1] == 0 {
			continue
		}
		out = append(out, (closes[i]-closes[i-1])/closes[i-1])
	}
	return out
}

// StdDev returns the population standard deviation of values. Returns 0 if len < 2.
func StdDev(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))
	var sqSum float64
	for _, v := range values {
		d := v - mean
		sqSum += d * d
	}
	return math.Sqrt(sqSum / float64(len(values)))
}

// RealisedVolatility returns the standard deviation of daily returns (closes). Daily, not annualized.
func RealisedVolatility(closes []float64) float64 {
	return StdDev(DailyReturns(closes))
}

// MaxDrawdown returns the maximum drawdown (0..1) over the close series: max of (peak - trough) / peak.
// peak is the running maximum of closes; trough is the close at the same or later index.
func MaxDrawdown(closes []float64) float64 {
	if len(closes) < 2 {
		return 0
	}
	peak := closes[0]
	maxDD := 0.0
	for i := 1; i < len(closes); i++ {
		if closes[i] > peak {
			peak = closes[i]
		}
		if peak <= 0 {
			continue
		}
		dd := (peak - closes[i]) / peak
		if dd > maxDD {
			maxDD = dd
		}
	}
	return maxDD
}
