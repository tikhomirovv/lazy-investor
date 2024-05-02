package logging

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

const (
	LogLevelDebug   = "debug"
	LogLevelInfo    = "info"
	LogLevelWarn    = "warn"
	LogLevelError   = "error"
	DefaultLogLevel = zerolog.DebugLevel
)

type ZLogger struct {
	logger zerolog.Logger
}

func NewLogger() *ZLogger {
	logger := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})).
		With().
		Timestamp().
		Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return &ZLogger{logger}
}

func (z ZLogger) Debug(msg string, args ...interface{}) {
	z.logger.Debug().Timestamp().Fields(args).Msg(msg)
}

func (z ZLogger) Info(msg string, args ...interface{}) {
	z.logger.Info().Timestamp().Fields(args).Msg(msg)
}

func (z ZLogger) Infof(msg string, args ...interface{}) {
	z.logger.Info().Timestamp().Fields(args).Msg(msg)
}

func (z ZLogger) Warn(msg string, args ...interface{}) {
	z.logger.Warn().Timestamp().Fields(args).Msg(msg)
}

func (z ZLogger) Error(msg string, args ...interface{}) {
	z.logger.Error().Timestamp().Fields(args).Msg(msg)
}

func (z ZLogger) Errorf(msg string, args ...interface{}) {
	z.logger.Error().Timestamp().Fields(args).Msg(msg)
}

func (z ZLogger) Panic(msg string, args ...interface{}) {
	z.logger.Panic().Timestamp().Fields(args).Msg(msg)
}

func (z ZLogger) Printf(format string, v ...interface{}) {
	z.logger.Printf(format, v...)
}

func (z ZLogger) Fatal(v ...interface{}) {
	z.logger.Fatal().Timestamp().Msg(fmt.Sprint(v...))
}

func (z ZLogger) Fatalf(format string, args ...interface{}) {
	z.logger.Fatal().Timestamp().Msgf(format, args...)
}

func (z ZLogger) Println(args ...interface{}) {
	z.logger.Info().Timestamp().Msgf("%v\r\n", args...)
}

func (z ZLogger) Print(args ...interface{}) {
	z.logger.Print(args...)
}
