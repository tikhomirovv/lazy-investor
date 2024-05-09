package chart

import (
	"io"
	"math/rand"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/dto"
	gc "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
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
}

// https://github.com/wcharczuk/go-chart/blob/main/examples/stock_analysis/main.go
func (cs *ChartService) Generate(chart *ChartValues, w io.Writer) error {
	var dates []time.Time
	var values []float64
	for _, c := range chart.Candles {
		dates = append(dates, c.Time)
		values = append(values, c.Close)
	}

	var datesTrends []time.Time
	var valuesTrends []float64
	for _, t := range chart.Trends {
		datesTrends = append(datesTrends, t.Candle.Time)
		valuesTrends = append(valuesTrends, t.Candle.Close)
	}
	priceSeries := gc.TimeSeries{
		Name: "SPY",
		Style: gc.Style{
			Show:        true,
			StrokeColor: gc.GetDefaultColor(0),
		},
		XValues: dates,
		YValues: values,
	}

	trendSeries := gc.TimeSeries{
		Name: "Trends",
		Style: gc.Style{
			StrokeWidth: gc.Disabled,
			DotWidth:    5,
			Show:        true,
			// StrokeColor: gc.GetDefaultColor(1),
		},
		XValues: datesTrends,
		YValues: valuesTrends,
	}
	smaSeries := gc.SMASeries{
		Name: "SPY - SMA",
		Style: gc.Style{
			Show:            true,
			StrokeColor:     drawing.ColorRed,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: priceSeries,
	}
	bbSeries := &gc.BollingerBandsSeries{
		Name: "SPY - Bol. Bands",
		Style: gc.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex("efefef"),
			FillColor:   drawing.ColorFromHex("efefef").WithAlpha(64),
		},
		InnerSeries: priceSeries,
	}

	series := []gc.Series{bbSeries,
		priceSeries,
		smaSeries,
		trendSeries,
	}

	for _, ema := range chart.EMAs {
		series = append(series, getEMATimeSeries(ema.Name, ema.Dates, ema.Values))
	}

	min, max := findMinMax(values)
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

func getEMATimeSeries(name string, dates []time.Time, values []float64) gc.TimeSeries {
	colorIndex := rand.Intn(4)
	return gc.TimeSeries{
		Name: name,
		Style: gc.Style{
			Show:        true,
			StrokeColor: gc.GetDefaultColor(colorIndex),
		},
		XValues: dates,
		YValues: values,
	}
}
