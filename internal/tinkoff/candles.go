package tinkoff

import (
	"github.com/tikhomirovv/lazy-investor/internal/dto"
	pb "github.com/tinkoff/invest-api-go-sdk/proto"
)

// Интервал свечей.
type CandleInterval int32

const (
	CandleIntervalUnspecified CandleInterval = 0  //Интервал не определён.
	CandleInterval1Min        CandleInterval = 1  //1 минута.
	CandleInterval5Min        CandleInterval = 2  //5 минут.
	CandleInterval15Min       CandleInterval = 3  //15 минут.
	CandleIntervalHour        CandleInterval = 4  //1 час.
	CandleIntervalDay         CandleInterval = 5  //1 день.
	CandleInterval2Min        CandleInterval = 6  //2 минуты.
	CandleInterval3Min        CandleInterval = 7  //3 минуты.
	CandleInterval10Min       CandleInterval = 8  //10 минут.
	CandleInterval30Min       CandleInterval = 9  //30 минут.
	CandleInterval2Hour       CandleInterval = 10 //2 часа.
	CandleInterval4Hour       CandleInterval = 11 //4 часа.
	CandleIntervalWeek        CandleInterval = 12 //1 неделя.
	CandleIntervalMonth       CandleInterval = 13 //1 месяц.
)

func Map(candles []*pb.HistoricCandle) []*dto.Candle {
	var result []*dto.Candle
	for _, c := range candles {
		result = append(result, &dto.Candle{
			Open:       c.Open.ToFloat(),
			High:       c.High.ToFloat(),
			Low:        c.Low.ToFloat(),
			Close:      c.Close.ToFloat(),
			Volume:     c.GetVolume(),
			Time:       c.GetTime().AsTime(),
			IsComplete: c.GetIsComplete(),
		})
	}
	return result
}
