package application

import (
	"github.com/tikhomirovv/lazy-investor/internal/tinkoff"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

type Application struct {
	logger  logging.Logger
	tinkoff *tinkoff.TinkoffService
}

func NewApplication(logger logging.Logger, tinkoff *tinkoff.TinkoffService) *Application {
	return &Application{
		logger:  logger,
		tinkoff: tinkoff,
	}
}

func (a *Application) Run() {
	a.tinkoff.Test()
}
