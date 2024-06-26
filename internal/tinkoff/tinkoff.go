package tinkoff

import (
	"context"
	"fmt"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/tikhomirovv/lazy-investor/internal/dto"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
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
func (t *TinkoffService) GetCandles(instrumentId string, from time.Time, to time.Time, interval CandleInterval) ([]dto.Candle, error) {
	marketDataService := t.client.NewMarketDataServiceClient()
	candlesResp, err := marketDataService.GetCandles(instrumentId, pb.CandleInterval(interval), from, to, pb.GetCandlesRequest_CANDLE_SOURCE_UNSPECIFIED)
	if err != nil {
		return nil, err
	}
	return Map(candlesResp.GetCandles()), nil
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
			t.logger.Debug("Instr", "i", i)
			return &dto.Instrument{
				Uid:  i.Uid,
				Name: i.Name,
				Isin: i.Isin,
			}, nil
		}
	}
	return nil, nil
}
