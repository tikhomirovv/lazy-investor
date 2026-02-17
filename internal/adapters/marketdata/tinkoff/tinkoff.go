// Package tinkoff implements market data and instrument lookup via Tinkoff Invest API.
// Moved from internal/services for SPEC.md adapters/marketdata/tinkoff layout.
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

// Config holds Tinkoff API connection settings (from env or config).
type Config struct {
	AppName string
	Host    string
	Token   string
}

// Service wraps invest-api-go-sdk client for candles and instruments.
type Service struct {
	config Config
	logger logging.Logger
	client *investgo.Client
}

// NewService creates Tinkoff SDK client and returns the adapter.
func NewService(config Config, logger logging.Logger) (*Service, error) {
	ctx := context.Background()
	client, err := investgo.NewClient(ctx, investgo.Config{
		AppName:  config.AppName,
		EndPoint: config.Host,
		Token:    config.Token,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("tinkoff client: %w", err)
	}
	return &Service{
		config: config,
		logger: logger,
		client: client,
	}, nil
}

// Stop closes the client connection (for graceful shutdown).
func (t *Service) Stop() {
	t.logger.Info("Closing client connection...")
	if err := t.client.Stop(); err != nil {
		t.logger.Error("Client shutdown error", "error", err)
	}
}

// GetCandles fetches historic candles for the instrument (default: day interval, last ~12 months).
func (t *Service) GetCandles(instrument *dto.Instrument) ([]dto.Candle, error) {
	var candles []dto.Candle
	dates := [][]time.Time{
		{
			time.Now().Add(-24 * 30 * 12 * time.Hour),
			time.Now().Add(-6 * time.Hour),
		},
	}
	for _, date := range dates {
		marketDataService := t.client.NewMarketDataServiceClient()
		candlesResp, err := marketDataService.GetCandles(instrument.Uid, pb.CandleInterval(CandleIntervalDay), date[0], date[1], pb.GetCandlesRequest_CANDLE_SOURCE_UNSPECIFIED)
		if err != nil {
			return nil, fmt.Errorf("GetCandles: %w", err)
		}
		candles = append(candles, MapCandles(candlesResp.GetCandles())...)
	}
	return candles, nil
}

// GetInstrumentByQuery finds a tradeable instrument by ISIN (or query string).
func (t *Service) GetInstrumentByQuery(q string) (*dto.Instrument, error) {
	instrumentService := t.client.NewInstrumentsServiceClient()
	resp, err := instrumentService.FindInstrument(q)
	if err != nil {
		return nil, err
	}
	for _, i := range resp.GetInstruments() {
		if i.ApiTradeAvailableFlag {
			return &dto.Instrument{
				Uid:  i.Uid,
				Name: i.Name,
				Isin: dto.Isin(i.Isin),
			}, nil
		}
	}
	return nil, nil
}

// CandleInterval matches Tinkoff proto enum for candle size.
type CandleInterval int32

const (
	CandleIntervalUnspecified CandleInterval = 0
	CandleInterval1Min        CandleInterval = 1
	CandleInterval5Min        CandleInterval = 2
	CandleInterval15Min       CandleInterval = 3
	CandleIntervalHour        CandleInterval = 4
	CandleIntervalDay         CandleInterval = 5
	CandleInterval2Min        CandleInterval = 6
	CandleInterval3Min        CandleInterval = 7
	CandleInterval10Min       CandleInterval = 8
	CandleInterval30Min       CandleInterval = 9
	CandleInterval2Hour       CandleInterval = 10
	CandleInterval4Hour       CandleInterval = 11
	CandleIntervalWeek        CandleInterval = 12
	CandleIntervalMonth       CandleInterval = 13
)

// MapCandles converts proto historic candles to domain DTOs.
func MapCandles(candles []*pb.HistoricCandle) []dto.Candle {
	result := make([]dto.Candle, 0, len(candles))
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
