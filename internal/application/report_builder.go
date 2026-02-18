// Package application provides report data and destination-specific formatters.
// Raw data is in ReportData; FormatForLog and FormatForTelegram produce strings for each destination.

package application

import (
	"fmt"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/application/features"
)

// InstrumentMetrics holds precomputed metrics and optional indicator/feature values for one instrument.
type InstrumentMetrics struct {
	Name        string
	Last        float64
	Change1d    float64
	Change7d    float64
	Change30d   float64
	Min         float64
	Max         float64
	AvgVolume   float64
	Volatility  float64
	MaxDrawdown float64
	// From IndicatorProvider (0 = not available); kept for backward compatibility
	SMA20 float64
	EMA20 float64
	RSI14 float64
	// EMA feature set (EMA20, EMA100, ema20_above_ema100); nil if not computed
	EMAFeatures *features.EMAFeatureSet
}

// ReportData is the raw report input for formatters (log vs Telegram).
type ReportData struct {
	AsOf  time.Time
	Rows  []InstrumentMetrics
}

// FormatForLog returns a report string for debug/log output (verbose, multiline).
func FormatForLog(data ReportData) string {
	const dateFmt = "2006-01-02 15:04"
	out := fmt.Sprintf("Stage 0 Report — %s\n\n", data.AsOf.Format(dateFmt))
	for _, r := range data.Rows {
		out += fmt.Sprintf("** %s\n", r.Name)
		out += fmt.Sprintf("   Last: %.2f  |  1d: %+.2f%%  7d: %+.2f%%  30d: %+.2f%%\n", r.Last, r.Change1d, r.Change7d, r.Change30d)
		out += fmt.Sprintf("   Range: %.2f — %.2f  |  AvgVol: %.0f  |  Vol(day): %.4f  MaxDD: %.1f%%\n", r.Min, r.Max, r.AvgVolume, r.Volatility, r.MaxDrawdown*100)
		if r.SMA20 != 0 || r.EMA20 != 0 || r.RSI14 != 0 {
			out += fmt.Sprintf("   SMA20: %.2f  EMA20: %.2f  RSI14: %.1f\n", r.SMA20, r.EMA20, r.RSI14)
		}
		formatEMAFeaturesLog(&r, &out)
		out += "\n"
	}
	return out
}

// formatEMAFeaturesLog appends EMA feature lines (EMA20, EMA100, trend filter, events) to out.
func formatEMAFeaturesLog(r *InstrumentMetrics, out *string) {
	if r.EMAFeatures == nil {
		return
	}
	ef := r.EMAFeatures
	if ef.EMA20.Ready {
		*out += fmt.Sprintf("   EMA20: %.2f  Δ %.2f", ef.EMA20.Value, ef.EMA20.Delta)
		if len(ef.EMA20.Events) > 0 {
			*out += fmt.Sprintf("  [%v]", ef.EMA20.Events)
		}
		*out += "\n"
	}
	if ef.EMA100.Ready {
		*out += fmt.Sprintf("   EMA100: %.2f  Δ %.2f", ef.EMA100.Value, ef.EMA100.Delta)
		if len(ef.EMA100.Events) > 0 {
			*out += fmt.Sprintf("  [%v]", ef.EMA100.Events)
		}
		*out += "\n"
	}
	if ef.EMA20AboveReady {
		trend := "below"
		if ef.EMA20AboveEMA100 {
			trend = "above"
		}
		*out += fmt.Sprintf("   EMA20 %s EMA100", trend)
		if len(ef.CombinedEvents) > 0 {
			*out += fmt.Sprintf("  [%v]", ef.CombinedEvents)
		}
		*out += "\n"
	}
}

// FormatForTelegram returns a report string for Telegram (compact, one-line per instrument optional).
func FormatForTelegram(data ReportData) string {
	const dateFmt = "2006-01-02 15:04"
	out := fmt.Sprintf("📊 Report %s\n\n", data.AsOf.Format(dateFmt))
	for _, r := range data.Rows {
		out += fmt.Sprintf("• %s\n  %.2f  1d:%+.1f%% 7d:%+.1f%% 30d:%+.1f%%\n", r.Name, r.Last, r.Change1d, r.Change7d, r.Change30d)
		out += fmt.Sprintf("  Range %.2f–%.2f  Vol %.0f  MaxDD %.1f%%\n", r.Min, r.Max, r.AvgVolume, r.MaxDrawdown*100)
		if r.SMA20 != 0 || r.EMA20 != 0 || r.RSI14 != 0 {
			out += fmt.Sprintf("  SMA20 %.2f  EMA20 %.2f  RSI %.0f\n", r.SMA20, r.EMA20, r.RSI14)
		}
		formatEMAFeaturesTelegram(&r, &out)
		out += "\n"
	}
	return out
}

// formatEMAFeaturesTelegram appends compact EMA feature line for Telegram (EMA20, EMA100, trend, events).
func formatEMAFeaturesTelegram(r *InstrumentMetrics, out *string) {
	if r.EMAFeatures == nil {
		return
	}
	ef := r.EMAFeatures
	if ef.EMA20.Ready {
		*out += fmt.Sprintf("  EMA20 %.2f Δ%.2f", ef.EMA20.Value, ef.EMA20.Delta)
		if len(ef.EMA20.Events) > 0 {
			*out += fmt.Sprintf(" %v", ef.EMA20.Events)
		}
		*out += "\n"
	}
	if ef.EMA100.Ready {
		*out += fmt.Sprintf("  EMA100 %.2f Δ%.2f", ef.EMA100.Value, ef.EMA100.Delta)
		if len(ef.EMA100.Events) > 0 {
			*out += fmt.Sprintf(" %v", ef.EMA100.Events)
		}
		*out += "\n"
	}
	if ef.EMA20AboveReady {
		trend := "below"
		if ef.EMA20AboveEMA100 {
			trend = "above"
		}
		*out += fmt.Sprintf("  EMA20 %s EMA100", trend)
		if len(ef.CombinedEvents) > 0 {
			*out += fmt.Sprintf(" %v", ef.CombinedEvents)
		}
		*out += "\n"
	}
}
