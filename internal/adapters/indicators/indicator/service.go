// Package indicator implements ports.IndicatorProvider using github.com/cinar/indicator/v2.
// Only this package may import cinar/indicator; app and ports stay dependency-free.
package indicator

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/momentum"
	"github.com/cinar/indicator/v2/trend"
	"github.com/tikhomirovv/lazy-investor/internal/ports"
)

const (
	smaPeriod = 20
	emaPeriod = 20
	rsiPeriod = 14
)

// Ensure Service implements ports.IndicatorProvider.
var _ ports.IndicatorProvider = (*Service)(nil)

// Service computes SMA, EMA, RSI via cinar/indicator/v2.
type Service struct{}

// NewService creates the indicator adapter.
func NewService() *Service {
	return &Service{}
}

// Compute returns the last value of SMA(20), EMA(20), RSI(14). Short series may yield 0.
func (s *Service) Compute(closes []float64) ports.IndicatorValues {
	out := ports.IndicatorValues{}
	if len(closes) < smaPeriod {
		return out
	}
	// Library uses channels; convert slice to channel, compute, then take last value.
	smaSeries := helper.ChanToSlice(trend.NewSmaWithPeriod[float64](smaPeriod).Compute(helper.SliceToChan(closes)))
	if len(smaSeries) > 0 {
		out.SMA20 = smaSeries[len(smaSeries)-1]
	}
	if len(closes) < emaPeriod {
		return out
	}
	emaSeries := helper.ChanToSlice(trend.NewEmaWithPeriod[float64](emaPeriod).Compute(helper.SliceToChan(closes)))
	if len(emaSeries) > 0 {
		out.EMA20 = emaSeries[len(emaSeries)-1]
	}
	if len(closes) < rsiPeriod+1 {
		return out
	}
	rsiSeries := helper.ChanToSlice(momentum.NewRsiWithPeriod[float64](rsiPeriod).Compute(helper.SliceToChan(closes)))
	if len(rsiSeries) > 0 {
		out.RSI14 = rsiSeries[len(rsiSeries)-1]
	}
	return out
}

// EMA returns the full EMA series for the given period, aligned to closes by index.
// Values slice has length len(closes); indices 0..Warmup-1 are zero. Ready is true when at least two values exist.
func (s *Service) EMA(closes []float64, period int) ports.SeriesResult {
	out := ports.SeriesResult{Warmup: period, Values: make([]float64, len(closes))}
	if period <= 0 || len(closes) < period {
		return out
	}
	emaSeries := helper.ChanToSlice(trend.NewEmaWithPeriod[float64](period).Compute(helper.SliceToChan(closes)))
	// cinar/indicator typically returns series aligned from start; length may equal len(closes) or len(closes)-period+1.
	// We align so that Values[i] corresponds to closes[i]: pad front with zeros if needed.
	n := len(emaSeries)
	if n == 0 {
		return out
	}
	offset := len(closes) - n
	if offset < 0 {
		offset = 0
		n = len(closes)
	}
	for i := 0; i < n; i++ {
		out.Values[offset+i] = emaSeries[i]
	}
	out.Ready = len(closes) >= period && n >= 2
	return out
}
