package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lazy-investor/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/tinkoff/invest-api-go-sdk/investgo"
	pb "github.com/tinkoff/invest-api-go-sdk/proto"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	appName := os.Getenv("APP_NAME")
	host := os.Getenv("TINKOFF_API_HOST")
	token := os.Getenv("TINKOFF_API_TOKEN")
	config := investgo.Config{AppName: appName, EndPoint: host, Token: token}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	logger := &logger.Logger{}

	// создаем клиента для investAPI, он позволяет создавать нужные сервисы и уже
	// через них вызывать нужные методы
	client, err := investgo.NewClient(ctx, config, logger)
	if err != nil {
		logger.Errorf("client creating error %v", err.Error())
	}
	defer func() {
		logger.Infof("closing client connection")
		err := client.Stop()
		if err != nil {
			logger.Errorf("client shutdown error %v", err.Error())
		}
	}()

	// Разово получить котировки по инструменту
	// создаем клиента для сервиса маркетдаты
	MarketDataService := client.NewMarketDataServiceClient()
	from := time.Now().Add(-6 * time.Hour)
	to := time.Now()
	instrumentId := "BBG004730N88" // SBER
	candlesResp, err := MarketDataService.GetCandles(instrumentId, pb.CandleInterval_CANDLE_INTERVAL_15_MIN, from, to)
	if err != nil {
		logger.Errorf(err.Error())
	} else {
		candles := candlesResp.GetCandles()
		for i, candle := range candles {
			fmt.Printf("candle number %d, high: %v, open: %v, close: %v, low:%v, volume: %v\n", i, candle.GetHigh().ToFloat(), candle.GetOpen().ToFloat(), candle.GetClose().ToFloat(), candle.GetLow().ToFloat(), candle.GetVolume())
		}
	}

}
