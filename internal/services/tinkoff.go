package services

import (
	"context"
	"fmt"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/tikhomirovv/lazy-investor/internal/dto"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

type TinkoffConfig struct {
	AppName string
	Host    string
	Token   string
}

type TinkoffService struct {
	config TinkoffConfig
	logger logging.Logger
	client *investgo.Client
}

func NewTinkoffService(config TinkoffConfig, logger logging.Logger) (*TinkoffService, error) {
	// создаем клиента для investAPI, он позволяет создавать нужные сервисы и уже
	// через них вызывать нужные методы
	ctx := context.Background()
	client, err := investgo.NewClient(ctx, investgo.Config{
		AppName:  config.AppName,
		EndPoint: config.Host,
		Token:    config.Token,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("client creating error %w", err)
	}
	return &TinkoffService{
		config: config,
		logger: logger,
		client: client,
	}, nil
}

func (t *TinkoffService) Stop() {
	t.logger.Info("Closing client connection...")
	err := t.client.Stop()
	if err != nil {
		t.logger.Error("Client shutdown error", "error", err)
	}
}

// Разово получить котировки по инструменту
func (t *TinkoffService) GetCandles(instrument *dto.Instrument) ([]dto.Candle, error) {
	var candles []dto.Candle
	dates := [][]time.Time{
		// {
		// 	time.Now().Add(-24 * 30 * 24 * time.Hour),
		// 	time.Now().Add(-24 * 30 * 12 * time.Hour),
		// },
		{
			time.Now().Add(-24 * 30 * 12 * time.Hour),
			time.Now().Add(-6 * time.Hour),
		}}
	for _, date := range dates {
		marketDataService := t.client.NewMarketDataServiceClient()
		candlesResp, err := marketDataService.GetCandles(instrument.Uid, pb.CandleInterval(CandleIntervalDay), date[0], date[1], pb.GetCandlesRequest_CANDLE_SOURCE_UNSPECIFIED)
		if err != nil {
			return nil, fmt.Errorf("TinkoffService.getCandles: %w", err)
		}
		ccc := Map(candlesResp.GetCandles())
		if err != nil {
			t.logger.Error("GetCandles error", "error", err)
			return nil, fmt.Errorf("TinkoffService.getCandles: %w", err)
		}
		// t.logger.Debug("GetCandles", "candles", ccc)
		candles = append(candles, ccc...)
	}
	return candles, nil
}

func (t *TinkoffService) GetInstrumentIdByQuery(q string) (*dto.Instrument, error) {
	instrumentService := t.client.NewInstrumentsServiceClient()
	resp, err := instrumentService.FindInstrument(q)
	if err != nil {
		return nil, err
	}
	instruments := resp.GetInstruments()
	for _, i := range instruments {
		if i.ApiTradeAvailableFlag {
			// t.logger.Debug("Instr", "i", i)
			return &dto.Instrument{
				Uid:  i.Uid,
				Name: i.Name,
				Isin: dto.Isin(i.Isin),
			}, nil
		}
	}
	return nil, nil
}

// Candles
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

func Map(candles []*pb.HistoricCandle) []dto.Candle {
	var result []dto.Candle
	for _, c := range candles {
		result = append(result, dto.Candle{
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
