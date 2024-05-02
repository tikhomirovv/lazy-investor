package tinkoff

import (
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

type TinkoffService struct {
	logger logging.Logger
}

func NewTinkoffService(logger logging.Logger) *TinkoffService {
	return &TinkoffService{
		logger: logger,
	}
}

func (t *TinkoffService) Test() {
	t.logger.Debug("Tinkoff text")
	// fmt.Println("Tinkoff test")
}
