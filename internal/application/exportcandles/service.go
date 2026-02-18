// Service fetches candles via MarketDataProvider and returns CSV bytes.
// No knowledge of CLI or Telegram; callers pass instrument query, timeframe, and period.

package exportcandles

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/dto"
	"github.com/tikhomirovv/lazy-investor/internal/ports"
)

// Service exports candles to CSV. Depends only on MarketDataProvider.
type Service struct {
	market ports.MarketDataProvider
}

// NewService creates the export service. market must not be nil.
func NewService(market ports.MarketDataProvider) *Service {
	return &Service{market: market}
}

// Export fetches candles for the given instrument (ISIN or ticker), timeframe string, and [from, to),
// then returns CSV bytes with header: time,open,high,low,close,volume.
// Time is formatted as RFC3339.
func (s *Service) Export(ctx context.Context, instrumentQuery, timeframe string, from, to time.Time) ([]byte, error) {
	interval, err := ParseTimeframe(timeframe)
	if err != nil {
		return nil, err
	}
	if interval == ports.IntervalUnspecified {
		return nil, fmt.Errorf("timeframe %q is not supported", timeframe)
	}

	instrument, err := s.market.FindInstrument(ctx, instrumentQuery)
	if err != nil {
		return nil, fmt.Errorf("find instrument %q: %w", instrumentQuery, err)
	}
	if instrument == nil {
		return nil, fmt.Errorf("instrument not found: %q", instrumentQuery)
	}

	candles, err := s.market.GetCandles(ctx, instrument, from, to, interval)
	if err != nil {
		return nil, fmt.Errorf("get candles: %w", err)
	}
	if len(candles) == 0 {
		return nil, fmt.Errorf("no candles for %q in the given period", instrumentQuery)
	}

	return candlesToCSV(candles), nil
}

// candlesToCSV writes candles to CSV with header time,open,high,low,close,volume. RFC3339 for time.
func candlesToCSV(candles []dto.Candle) []byte {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Header row.
	_ = w.Write([]string{"time", "open", "high", "low", "close", "volume"})

	for _, c := range candles {
		row := []string{
			c.Time.UTC().Format(time.RFC3339),
			floatToStr(c.Open),
			floatToStr(c.High),
			floatToStr(c.Low),
			floatToStr(c.Close),
			fmt.Sprintf("%d", c.Volume),
		}
		_ = w.Write(row)
	}

	w.Flush()
	return buf.Bytes()
}

func floatToStr(f float64) string {
	return fmt.Sprintf("%.4f", f)
}
