package tinkoff

import (
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
}

func NewTinkoffService(config Config, logger logging.Logger) *TinkoffService {
	return &TinkoffService{
		config: config,
		logger: logger,
	}
}

func (t *TinkoffService) Test() {
	t.logger.Debug("Tinkoff test", "app", t.config.AppName)
	// fmt.Println("Tinkoff test")
}

func te() {

	// config := investgo.Config{AppName: appName, EndPoint: host, Token: token}

	// ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	// defer cancel()

	// logger := &logger.Logger{}

	// // создаем клиента для investAPI, он позволяет создавать нужные сервисы и уже
	// // через них вызывать нужные методы
	// client, err := investgo.NewClient(ctx, config, logger)
	// if err != nil {
	// 	logger.Errorf("client creating error %v", err.Error())
	// }
	// defer func() {
	// 	logger.Infof("closing client connection")
	// 	err := client.Stop()
	// 	if err != nil {
	// 		logger.Errorf("client shutdown error %v", err.Error())
	// 	}
	// }()

	// // Разово получить котировки по инструменту
	// // создаем клиента для сервиса маркетдаты
	// MarketDataService := client.NewMarketDataServiceClient()
	// from := time.Now().Add(-6 * time.Hour)
	// to := time.Now()
	// instrumentId := "BBG004730N88" // SBER
	// candlesResp, err := MarketDataService.GetCandles(instrumentId, pb.CandleInterval_CANDLE_INTERVAL_15_MIN, from, to)
	// if err != nil {
	// 	logger.Errorf(err.Error())
	// } else {
	// 	candles := candlesResp.GetCandles()
	// 	for i, candle := range candles {
	// 		fmt.Printf("candle number %d, high: %v, open: %v, close: %v, low:%v, volume: %v\n", i, candle.GetHigh().ToFloat(), candle.GetOpen().ToFloat(), candle.GetClose().ToFloat(), candle.GetLow().ToFloat(), candle.GetVolume())
	// 	}
	// }

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
