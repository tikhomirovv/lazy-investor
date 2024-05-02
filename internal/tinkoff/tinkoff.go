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
