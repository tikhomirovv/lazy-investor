package tinkoff

import (
	"context"
	"fmt"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/dto"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
	"github.com/tinkoff/invest-api-go-sdk/investgo"
	pb "github.com/tinkoff/invest-api-go-sdk/proto"
)

type Config struct {
	AppName string
	Host    string
	Token   string
}

type TinkoffService struct {
	config Config
	logger logging.Logger
	client *investgo.Client
}

func NewTinkoffService(config Config, logger logging.Logger) (*TinkoffService, error) {
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
func (t *TinkoffService) GetCandles(instrumentId string, from time.Time, to time.Time, interval CandleInterval) ([]*dto.Candle, error) {
	marketDataService := t.client.NewMarketDataServiceClient()
	candlesResp, err := marketDataService.GetCandles(instrumentId, pb.CandleInterval(interval), from, to)
	if err != nil {
		return nil, err
	}
	return Map(candlesResp.GetCandles()), nil
}

func te() {

	// // минутные свечи TCSG за последние двое суток
	// candles, err := MarketDataService.GetHistoricCandles(&investgo.GetHistoricCandlesRequest{
	// 	Instrument: instrumentId,
	// 	Interval:   pb.CandleInterval_CANDLE_INTERVAL_1_MIN,
	// 	From:       time.Date(2023, time.June, 2, 10, 0, 0, 0, time.UTC),
	// 	To:         time.Date(2023, time.June, 4, 0, 0, 0, 0, time.UTC),
	// 	File:       true,
	// 	FileName:   "sber_june_2_2023",
	// })
	// if err != nil {
	// 	logger.Errorf(err.Error())
	// } else {
	// 	for i, candle := range candles {
	// 		fmt.Printf("candle %v open = %v\n", i, candle.GetOpen().ToFloat())
	// 	}
	// }
}
