// Package tinkoff implements ports.MarketDataProvider via Tinkoff Invest API.
package tinkoff

import (
	"context"
	"fmt"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/tikhomirovv/lazy-investor/internal/dto"
	"github.com/tikhomirovv/lazy-investor/internal/ports"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

// Ensure Service implements ports.MarketDataProvider.
var _ ports.MarketDataProvider = (*Service)(nil)

// Config holds Tinkoff API connection settings (from env or config).
type Config struct {
	AppName string
	Host    string
	Token   string
}

// Service wraps invest-api-go-sdk client; implements MarketDataProvider.
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

// GetCandles fetches historic candles for the instrument in [from, to) with given interval.
func (t *Service) GetCandles(ctx context.Context, instrument *dto.Instrument, from, to time.Time, interval ports.CandleInterval) ([]dto.Candle, error) {
	marketDataService := t.client.NewMarketDataServiceClient()
	candlesResp, err := marketDataService.GetCandles(
		instrument.Uid,
		pb.CandleInterval(toProtoInterval(interval)),
		from,
		to,
		pb.GetCandlesRequest_CANDLE_SOURCE_UNSPECIFIED,
	)
	if err != nil {
		return nil, fmt.Errorf("GetCandles: %w", err)
	}
	return MapCandles(candlesResp.GetCandles()), nil
}

// FindInstrument finds a tradeable instrument by ISIN or query string.
func (t *Service) FindInstrument(ctx context.Context, query string) (*dto.Instrument, error) {
	instrumentService := t.client.NewInstrumentsServiceClient()
	resp, err := instrumentService.FindInstrument(query)
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

// toProtoInterval maps ports.CandleInterval to Tinkoff proto enum.
func toProtoInterval(interval ports.CandleInterval) int32 {
	return int32(interval)
}

// MapCandles converts proto historic candles to dto (contract between API and app).
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
