// export-candles is a standalone CLI that exports candles to CSV.
// Usage: instrument (ISIN or ticker), timeframe (e.g. 1d), and period (--from/--to or --last=N).
// Output: stdout or --output=path. Uses same MarketDataProvider (Tinkoff) as the main app.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tikhomirovv/lazy-investor/internal/application/exportcandles"
	"github.com/tikhomirovv/lazy-investor/pkg/wire"
)

const dateLayout = "2006-01-02"

func main() {
	_ = godotenv.Load()

	instrument := flag.String("instrument", "", "Instrument: ISIN or ticker (e.g. SBER, RU0009029540)")
	timeframe := flag.String("timeframe", "1d", "Candle timeframe: 1m, 5m, 15m, 1h, 1d, 1w, 1month")
	fromStr := flag.String("from", "", "Period start date YYYY-MM-DD (use with --to)")
	toStr := flag.String("to", "", "Period end date YYYY-MM-DD (use with --from)")
	lastDays := flag.Int("last", 0, "Last N days (alternative to --from/--to)")
	outputPath := flag.String("output", "", "Write CSV to file (default: stdout)")
	flag.Parse()

	if *instrument == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: -instrument is required")
		flag.Usage()
		os.Exit(1)
	}

	from, to, err := parsePeriod(*fromStr, *toStr, *lastDays)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	logger := wire.InitLogger()
	market, err := wire.InitTinkoffService(logger)
	if err != nil {
		logger.Errorf("init market data: %v", err)
		os.Exit(1)
	}
	defer market.Stop()

	svc := exportcandles.NewService(market)
	ctx := context.Background()

	csvBytes, err := svc.Export(ctx, *instrument, *timeframe, from, to)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *outputPath != "" {
		if err := os.WriteFile(*outputPath, csvBytes, 0644); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
			os.Exit(1)
		}
		logger.Info("wrote CSV", "path", *outputPath, "bytes", len(csvBytes))
		return
	}

	_, _ = os.Stdout.Write(csvBytes)
}

// parsePeriod returns from, to. Prefer --from/--to; if both empty use --last (default 30 days).
func parsePeriod(fromStr, toStr string, lastDays int) (from, to time.Time, err error) {
	if fromStr != "" && toStr != "" {
		from, err = time.Parse(dateLayout, fromStr)
		if err != nil {
			return from, to, fmt.Errorf("parse --from: %w", err)
		}
		to, err = time.Parse(dateLayout, toStr)
		if err != nil {
			return from, to, fmt.Errorf("parse --to: %w", err)
		}
		if !to.After(from) {
			return from, to, fmt.Errorf("--to must be after --from")
		}
		return from, to, nil
	}
	if fromStr != "" || toStr != "" {
		return from, to, fmt.Errorf("use both --from and --to, or use --last")
	}
	if lastDays <= 0 {
		lastDays = 30
	}
	to = time.Now()
	from = to.AddDate(0, 0, -lastDays)
	return from, to, nil
}
