package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/tikhomirovv/lazy-investor/pkg/wire"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// Application
	logger := wire.InitLogger()
	application, err := wire.InitApplication()
	if err != nil {
		panic(err)
	}

	go func() {
		application.Run(ctx)
	}()

	// Gracefull shutdown
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT)
	logger.Info("Application started")
	<-stopSignal
	logger.Info("Shutting down gracefully...")
	cancel()
	application.Stop()
	time.Sleep(1 * time.Second)
	// Завершение работы
	logger.Info("Shutdown finished")
	os.Exit(0)
}
