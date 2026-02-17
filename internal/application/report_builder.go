// Package application provides the main app orchestration and Stage 0 report building.
// ReportBuilder builds plain-text report from precomputed per-instrument metrics.

package application

import (
	"fmt"
	"time"
)

// InstrumentMetrics holds precomputed metrics for one instrument (used by report builder).
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
}

// BuildReport produces a plain-text report from instrument metrics. Title includes asOf time.
func BuildReport(rows []InstrumentMetrics, asOf time.Time) string {
	const dateFmt = "2006-01-02 15:04"
	out := fmt.Sprintf("Stage 0 Report — %s\n\n", asOf.Format(dateFmt))
	for _, r := range rows {
		out += fmt.Sprintf("** %s\n", r.Name)
		out += fmt.Sprintf("   Last: %.2f  |  1d: %+.2f%%  7d: %+.2f%%  30d: %+.2f%%\n", r.Last, r.Change1d, r.Change7d, r.Change30d)
		out += fmt.Sprintf("   Range: %.2f — %.2f  |  AvgVol: %.0f  |  Vol(day): %.4f  MaxDD: %.1f%%\n\n", r.Min, r.Max, r.AvgVolume, r.Volatility, r.MaxDrawdown*100)
	}
	return out
}
