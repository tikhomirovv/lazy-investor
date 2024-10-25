package chart

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/dto"
	gc "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	util "github.com/wcharczuk/go-chart/util"
)

type ChartService struct{}

func NewChartService() *ChartService {
	return &ChartService{}
}

type ChartValues struct {
	Title   string
	Candles []dto.Candle
	Trends  []dto.TrendChange
	EMAs    []dto.EMA
	Swings  []dto.Swing
	ZigZags []dto.ZigZagPoint
}

// https://github.com/wcharczuk/go-chart/blob/main/examples/stock_analysis/main.go
func (cs *ChartService) Generate(chart *ChartValues, w io.Writer) error {

	close, high, low := getCandlesPrices(chart.Candles)

	// smaSeries := gc.SMASeries{
	// 	Name: "SPY - SMA",
	// 	Style: gc.Style{
	// 		Show:            true,
	// 		StrokeColor:     drawing.ColorRed,
	// 		StrokeDashArray: []float64{5.0, 5.0},
	// 	},
	// 	InnerSeries: close,
	// }
	bbSeries := &gc.BollingerBandsSeries{
		Name: "SPY - Bol. Bands",
		Style: gc.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex("efefef"),
			FillColor:   drawing.ColorFromHex("efefef").WithAlpha(64),
		},
		InnerSeries: close,
	}

	// swingsHigh, swingsLow := getSwingsTimeSeries(chart.Swings)
	zzHigh, zzLow := getZigZagsTimeSeries(chart.ZigZags)
	// trendsUp, trendsDown, trendsNo, _ := getTrendsTimeSeries(chart.Trends)
	series := []gc.Series{bbSeries,
		close, high, low,
		// smaSeries,
		// trendsUp, trendsDown, trendsNo,
		// swingsHigh, swingsLow,
		zzHigh, zzLow,
		// trendAnnotations,
	}

	for _, ema := range chart.EMAs {
		series = append(series, getEMATimeSeries(ema))
	}

	min, max := findMinMax(close.YValues)
	graph := gc.Chart{
		Title: chart.Title,
		TitleStyle: gc.Style{
			Show: true,
		},
		XAxis: gc.XAxis{
			Style: gc.Style{
				Show: true,
			},
			TickPosition: gc.TickPositionBetweenTicks,
		},
		YAxis: gc.YAxis{
			Style: gc.Style{
				Show: true,
			},
			Range: &gc.ContinuousRange{
				Max: max + 1,
				Min: min - 1,
			},
		},
		Series: series,
	}

	graph.Elements = []gc.Renderable{
		gc.Legend(&graph),
	}
	return graph.Render(gc.PNG, w)
}

func findMinMax(slice []float64) (min, max float64) {
	if len(slice) == 0 {
		return 0, 0
	}
	min, max = slice[0], slice[0]
	for _, value := range slice {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func getCandlesPrices(candles []dto.Candle) (close gc.TimeSeries, high gc.TimeSeries, low gc.TimeSeries) {
	var dates []time.Time
	var closeValues []float64
	var highValues []float64
	var lowValues []float64
	for _, c := range candles {
		dates = append(dates, c.Time)
		closeValues = append(closeValues, c.Close)
		highValues = append(highValues, c.High)
		lowValues = append(lowValues, c.Low)
	}
	close = gc.TimeSeries{
		Name: "Price Close",
		Style: gc.Style{
			Show:        true,
			StrokeColor: gc.GetDefaultColor(0),
		},
		XValues: dates,
		YValues: closeValues,
	}
	high = gc.TimeSeries{
		Name: "Price High",
		Style: gc.Style{
			Show:            true,
			StrokeDashArray: []float64{5.0, 5.0},
			StrokeColor:     gc.GetDefaultColor(1),
		},
		XValues: dates,
		YValues: highValues,
	}
	low = gc.TimeSeries{
		Name: "Price Low",
		Style: gc.Style{
			Show:            true,
			StrokeDashArray: []float64{5.0, 5.0},
			StrokeColor:     gc.GetDefaultColor(2),
		},
		XValues: dates,
		YValues: lowValues,
	}
	return
}

func getEMATimeSeries(ema dto.EMA) gc.TimeSeries {
	colorIndex := rand.Intn(4)
	return gc.TimeSeries{
		Name: ema.Name,
		Style: gc.Style{
			Show:        true,
			StrokeColor: gc.GetDefaultColor(colorIndex),
		},
		XValues: ema.Dates,
		YValues: ema.Values,
	}
}

func getTrendsTimeSeries(trends []dto.TrendChange) (up gc.TimeSeries, down gc.TimeSeries, no gc.TimeSeries, annotations gc.AnnotationSeries) {
	var upDates, downDates, noDates []time.Time
	var upValues, downValues, noValues []float64
	var annotationValues []gc.Value2
	for _, t := range trends {
		annotationValues = append(annotationValues, gc.Value2{
			XValue: util.Time.ToFloat64(t.Swing.Candle.Time),
			YValue: t.Swing.GetValue(),
			Label:  fmt.Sprintf("Trend: %s", t.Trend.String()),
		})
		if t.Trend == dto.TrendUp {
			upDates = append(upDates, t.Swing.Candle.Time)
			upValues = append(upValues, t.Swing.GetValue())
		} else if t.Trend == dto.TrendDown {
			downDates = append(downDates, t.Swing.Candle.Time)
			downValues = append(downValues, t.Swing.GetValue())
		} else {
			noDates = append(noDates, t.Swing.Candle.Time)
			noValues = append(noValues, t.Swing.GetValue())
		}
	}
	up = gc.TimeSeries{
		Name: "Trends Up",
		Style: gc.Style{
			StrokeWidth: gc.Disabled,
			DotWidth:    7,
			DotColor:    gc.GetAlternateColor(1),
			Show:        true,
		},
		XValues: upDates,
		YValues: upValues,
	}
	down = gc.TimeSeries{
		Name: "Trends Down",
		Style: gc.Style{
			StrokeWidth: gc.Disabled,
			DotWidth:    7,
			DotColor:    gc.GetDefaultColor(2),
			Show:        true,
		},
		XValues: downDates,
		YValues: downValues,
	}
	no = gc.TimeSeries{
		Name: "Trends No",
		Style: gc.Style{
			StrokeWidth: gc.Disabled,
			DotWidth:    7,
			DotColor:    gc.GetAlternateColor(3),
			Show:        true,
		},
		XValues: noDates,
		YValues: noValues,
	}
	annotations = gc.AnnotationSeries{
		Annotations: annotationValues,
	}
	return
}

func getSwingsTimeSeries(swings []dto.Swing) (gc.TimeSeries, gc.TimeSeries) {
	var hDates, lDates []time.Time
	var hValues, lValues []float64
	for _, s := range swings {
		if s.Type == dto.SwingHigh {
			hDates = append(hDates, s.Candle.Time)
			hValues = append(hValues, s.Candle.High)
		} else {
			lDates = append(lDates, s.Candle.Time)
			lValues = append(lValues, s.Candle.Low)
		}
	}
	return gc.TimeSeries{
			Name: "SwingsHigh",
			Style: gc.Style{
				StrokeWidth: gc.Disabled,
				DotWidth:    2,
				DotColor:    gc.GetDefaultColor(1),
				Show:        true,
			},
			XValues: hDates,
			YValues: hValues,
		},
		gc.TimeSeries{
			Name: "SwingsLow",
			Style: gc.Style{
				StrokeWidth: gc.Disabled,
				DotWidth:    2,
				DotColor:    gc.GetDefaultColor(2),
				Show:        true,
			},
			XValues: lDates,
			YValues: lValues,
		}
}

func getZigZagsTimeSeries(zz []dto.ZigZagPoint) (gc.TimeSeries, gc.TimeSeries) {
	var hDates, lDates []time.Time
	var hValues, lValues []float64
	for _, z := range zz {
		if z.Type == dto.ZigZagHigh {
			hDates = append(hDates, z.Candle.Time)
			hValues = append(hValues, z.Candle.High)
		} else {
			lDates = append(lDates, z.Candle.Time)
			lValues = append(lValues, z.Candle.Low)
		}
	}
	return gc.TimeSeries{
			Name: "ZZHigh",
			Style: gc.Style{
				StrokeWidth: gc.Disabled,
				DotWidth:    5,
				DotColor:    gc.GetDefaultColor(1),
				Show:        true,
			},
			XValues: hDates,
			YValues: hValues,
		},
		gc.TimeSeries{
			Name: "ZZLow",
			Style: gc.Style{
				StrokeWidth: gc.Disabled,
				DotWidth:    5,
				DotColor:    gc.GetDefaultColor(2),
				Show:        true,
			},
			XValues: lDates,
			YValues: lValues,
		}
}
