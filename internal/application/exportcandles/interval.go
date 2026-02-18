// Package exportcandles provides a service to fetch candles and export them as CSV.
// Used by both CLI (cmd/export-candles) and Telegram bot (/candles command).

package exportcandles

import (
	"fmt"
	"strings"

	"github.com/tikhomirovv/lazy-investor/internal/ports"
)

// ParseTimeframe maps a string (e.g. "1d", "1h") to ports.CandleInterval.
// Used by both CLI and Telegram handler so parsing lives in one place.
func ParseTimeframe(s string) (ports.CandleInterval, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1m":
		return ports.Interval1Min, nil
	case "2m":
		return ports.Interval2Min, nil
	case "3m":
		return ports.Interval3Min, nil
	case "5m":
		return ports.Interval5Min, nil
	case "10m":
		return ports.Interval10Min, nil
	case "15m":
		return ports.Interval15Min, nil
	case "30m":
		return ports.Interval30Min, nil
	case "1h":
		return ports.Interval1Hour, nil
	case "2h":
		return ports.Interval2Hour, nil
	case "4h":
		return ports.Interval4Hour, nil
	case "1d":
		return ports.Interval1Day, nil
	case "1w":
		return ports.Interval1Week, nil
	case "1month", "1mo":
		return ports.Interval1Month, nil
	default:
		return ports.IntervalUnspecified, fmt.Errorf("unknown timeframe %q (use 1m, 5m, 15m, 1h, 1d, 1w, 1month)", s)
	}
}
