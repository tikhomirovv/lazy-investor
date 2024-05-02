package tinkoff

import (
	"github.com/tikhomirovv/lazy-investor/internal/dto"
	pb "github.com/tinkoff/invest-api-go-sdk/proto"
)

// Интервал свечей.
type CandleInterval int32

const (
	CANDLE_INTERVAL_UNSPECIFIED CandleInterval = 0  //Интервал не определён.
	CANDLE_INTERVAL_1_MIN       CandleInterval = 1  //1 минута.
	CANDLE_INTERVAL_5_MIN       CandleInterval = 2  //5 минут.
	CANDLE_INTERVAL_15_MIN      CandleInterval = 3  //15 минут.
	CANDLE_INTERVAL_HOUR        CandleInterval = 4  //1 час.
	CANDLE_INTERVAL_DAY         CandleInterval = 5  //1 день.
	CANDLE_INTERVAL_2_MIN       CandleInterval = 6  //2 минуты.
	CANDLE_INTERVAL_3_MIN       CandleInterval = 7  //3 минуты.
	CANDLE_INTERVAL_10_MIN      CandleInterval = 8  //10 минут.
	CANDLE_INTERVAL_30_MIN      CandleInterval = 9  //30 минут.
	CANDLE_INTERVAL_2_HOUR      CandleInterval = 10 //2 часа.
	CANDLE_INTERVAL_4_HOUR      CandleInterval = 11 //4 часа.
	CANDLE_INTERVAL_WEEK        CandleInterval = 12 //1 неделя.
	CANDLE_INTERVAL_MONTH       CandleInterval = 13 //1 месяц.
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
