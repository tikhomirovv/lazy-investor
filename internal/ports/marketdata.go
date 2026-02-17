// Package ports defines interfaces (contracts) between app/domain and adapters.
// MarketDataProvider: get candles and find instrument; contract uses dto (layer between exchange API and app).

package ports

import (
	"context"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/dto"
)

// CandleInterval is exchange-agnostic; adapters map to broker-specific enums.
type CandleInterval int32

const (
	IntervalUnspecified CandleInterval = 0
	Interval1Min        CandleInterval = 1
	Interval5Min        CandleInterval = 2
	Interval15Min       CandleInterval = 3
	Interval1Hour       CandleInterval = 4
	Interval1Day        CandleInterval = 5
	Interval2Min        CandleInterval = 6
	Interval3Min        CandleInterval = 7
	Interval10Min       CandleInterval = 8
	Interval30Min       CandleInterval = 9
	Interval2Hour       CandleInterval = 10
	Interval4Hour       CandleInterval = 11
	Interval1Week       CandleInterval = 12
	Interval1Month      CandleInterval = 13
)

// MarketDataProvider fetches candles and instruments. Implemented by adapters (e.g. Tinkoff).
// DTOs (dto.Candle, dto.Instrument) are the contract between exchange API and application.
type MarketDataProvider interface {
	// GetCandles returns historic candles for the instrument in [from, to) with given interval.
	GetCandles(ctx context.Context, instrument *dto.Instrument, from, to time.Time, interval CandleInterval) ([]dto.Candle, error)
	// FindInstrument finds a tradeable instrument by ISIN or query string.
	FindInstrument(ctx context.Context, query string) (*dto.Instrument, error)
	// Stop closes the provider (e.g. broker connection). Safe to call multiple times.
	Stop()
}
