// Package chart implements PNG chart rendering only. No domain/dto types;
// caller passes drawing-oriented input (times, OHLC slices, optional overlays and point markers).
package chart

import (
	"io"
	"math/rand"
	"time"

	gc "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// Service renders price charts to PNG.
type Service struct{}

// NewService returns a chart renderer.
func NewService() *Service {
	return &Service{}
}

// LineSeries is a single line overlay (e.g. EMA): name + time/value pairs.
type LineSeries struct {
	Name   string
	Times  []time.Time
	Values []float64
}

// Point is a single marker (e.g. swing high/low, zigzag).
type Point struct {
	Time  time.Time
	Value float64
}

// Input describes what to draw: OHLC data plus optional Bollinger, line overlays, and point markers.
// All slice lengths for Times/Open/High/Low/Close must match. Overlays and points may be nil/empty.
type Input struct {
	Title       string
	Times       []time.Time
	Open        []float64
	High        []float64
	Low         []float64
	Close       []float64
	Bollinger   bool         // draw Bollinger Bands
	Overlays    []LineSeries // e.g. EMAs
	HighPoints  []Point      // optional scatter (e.g. zigzag high)
	LowPoints   []Point      // optional scatter (e.g. zigzag low)
}

// Generate writes a PNG chart to w. Only drawing logic; no dependency on internal/dto.
func (s *Service) Generate(in *Input, w io.Writer) error {
	if in == nil || len(in.Times) == 0 {
		return nil
	}
	closeSeries, highSeries, lowSeries := seriesFromOHLC(in.Times, in.Open, in.High, in.Low, in.Close)
	series := []gc.Series{closeSeries, highSeries, lowSeries}

	if in.Bollinger {
		bb := &gc.BollingerBandsSeries{
			Name: "Bollinger Bands",
			Style: gc.Style{
				Show:        true,
				StrokeColor: drawing.ColorFromHex("efefef"),
				FillColor:   drawing.ColorFromHex("efefef").WithAlpha(64),
			},
			InnerSeries: closeSeries,
		}
		series = append([]gc.Series{bb}, series...)
	}

	for _, ov := range in.Overlays {
		series = append(series, lineToTimeSeries(ov))
	}
	highPts, lowPts := pointsToTimeSeries(in.HighPoints, in.LowPoints)
	if len(in.HighPoints) > 0 {
		series = append(series, highPts)
	}
	if len(in.LowPoints) > 0 {
		series = append(series, lowPts)
	}

	min, max := findMinMax(closeSeries.YValues)
	graph := gc.Chart{
		Title:      in.Title,
		TitleStyle: gc.Style{Show: true},
		XAxis:      gc.XAxis{Style: gc.Style{Show: true}, TickPosition: gc.TickPositionBetweenTicks},
		YAxis: gc.YAxis{
			Style: gc.Style{Show: true},
			Range: &gc.ContinuousRange{Max: max + 1, Min: min - 1},
		},
		Series: series,
	}
	graph.Elements = []gc.Renderable{gc.Legend(&graph)}
	return graph.Render(gc.PNG, w)
}

func findMinMax(slice []float64) (min, max float64) {
	if len(slice) == 0 {
		return 0, 0
	}
	min, max = slice[0], slice[0]
	for _, v := range slice {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

func seriesFromOHLC(times []time.Time, open, high, low, close []float64) (closeSeries, highSeries, lowSeries gc.TimeSeries) {
	closeSeries = gc.TimeSeries{
		Name: "Close", Style: gc.Style{Show: true, StrokeColor: gc.GetDefaultColor(0)},
		XValues: times, YValues: close,
	}
	highSeries = gc.TimeSeries{
		Name: "High", Style: gc.Style{Show: true, StrokeDashArray: []float64{5, 5}, StrokeColor: gc.GetDefaultColor(1)},
		XValues: times, YValues: high,
	}
	lowSeries = gc.TimeSeries{
		Name: "Low", Style: gc.Style{Show: true, StrokeDashArray: []float64{5, 5}, StrokeColor: gc.GetDefaultColor(2)},
		XValues: times, YValues: low,
	}
	return closeSeries, highSeries, lowSeries
}

func lineToTimeSeries(l LineSeries) gc.TimeSeries {
	return gc.TimeSeries{
		Name: l.Name,
		Style: gc.Style{
			Show: true, StrokeColor: gc.GetDefaultColor(rand.Intn(4)),
		},
		XValues: l.Times,
		YValues: l.Values,
	}
}

func pointsToTimeSeries(high, low []Point) (highSeries, lowSeries gc.TimeSeries) {
	var hTimes, lTimes []time.Time
	var hV, lV []float64
	for _, p := range high {
		hTimes = append(hTimes, p.Time)
		hV = append(hV, p.Value)
	}
	for _, p := range low {
		lTimes = append(lTimes, p.Time)
		lV = append(lV, p.Value)
	}
	if len(hTimes) > 0 {
		highSeries = gc.TimeSeries{
			Name: "High points",
			Style: gc.Style{StrokeWidth: gc.Disabled, DotWidth: 5, DotColor: gc.GetDefaultColor(1), Show: true},
			XValues: hTimes, YValues: hV,
		}
	}
	if len(lTimes) > 0 {
		lowSeries = gc.TimeSeries{
			Name: "Low points",
			Style: gc.Style{StrokeWidth: gc.Disabled, DotWidth: 5, DotColor: gc.GetDefaultColor(2), Show: true},
			XValues: lTimes, YValues: lV,
		}
	}
	return highSeries, lowSeries
}
