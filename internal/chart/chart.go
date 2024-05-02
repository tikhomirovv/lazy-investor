package chart

import (
	"io"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/dto"
	gc "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

type ChartService struct{}

func NewChartService() *ChartService {
	return &ChartService{}
}

// https://github.com/wcharczuk/go-chart/blob/main/examples/stock_analysis/main.go
func (cs *ChartService) Generate(candles []dto.Candle, trends []dto.TrendChange, shortMA []float64, longMA []float64, w io.Writer) error {
	var dates []time.Time
	var values []float64
	var valuesShortMA []float64
	var valuesLongMA []float64
	for i, c := range candles {
		dates = append(dates, c.Time)
		values = append(values, c.Close)
		valuesShortMA = append(valuesShortMA, shortMA[i])
		valuesLongMA = append(valuesLongMA, longMA[i])
	}

	var datesTrends []time.Time
	var valuesTrends []float64
	for _, t := range trends {
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

	shortMASeries := gc.TimeSeries{
		Name: "ShortMA",
		Style: gc.Style{
			Show:        true,
			StrokeColor: gc.GetDefaultColor(1),
		},
		XValues: dates,
		YValues: valuesShortMA,
	}
	longMASeries := gc.TimeSeries{
		Name: "LongMA",
		Style: gc.Style{
			Show:        true,
			StrokeColor: gc.GetDefaultColor(2),
		},
		XValues: dates,
		YValues: valuesLongMA,
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

	min, max := findMinMax(values)
	graph := gc.Chart{
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
		Series: []gc.Series{
			bbSeries,
			priceSeries,
			smaSeries,
			trendSeries,
			shortMASeries,
			longMASeries,
		},
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
