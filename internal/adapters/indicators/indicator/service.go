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
