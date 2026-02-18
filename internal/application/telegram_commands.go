// telegram_commands.go handles incoming Telegram messages when handleCommands is enabled (e.g. /candles).

package application

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const candlesCommand = "/candles"

// runTelegramCommandListener starts long polling and dispatches messages to handleTelegramCommand.
func (a *Application) runTelegramCommandListener(ctx context.Context) {
	a.logger.Info("Telegram command listener started")
	a.telegram.ListenForMessages(ctx, a.handleTelegramCommand)
	a.logger.Info("Telegram command listener stopped")
}

// handleTelegramCommand processes one message. Supports /candles <instrument> <tf> <last_days>.
func (a *Application) handleTelegramCommand(chatID int64, text string) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, candlesCommand) {
		return
	}

	// Optional: restrict to allowed chat only.
	if a.config.Telegram.AllowedChatID != "" {
		allowed, err := strconv.ParseInt(a.config.Telegram.AllowedChatID, 10, 64)
		if err == nil && chatID != allowed {
			a.logger.Debug("ignoring command from non-allowed chat", "chatID", chatID)
			return
		}
	}

	args := strings.Fields(strings.TrimPrefix(text, candlesCommand))
	instrument, timeframe, lastDays, err := parseCandlesArgs(args)
	if err != nil {
		_ = a.telegram.SendMessageToChat(context.Background(), chatID, "Usage: /candles <instrument> <timeframe> <last_days>\nExample: /candles SBER 1d 30\nError: "+err.Error())
		return
	}

	to := time.Now()
	from := to.AddDate(0, 0, -lastDays)

	csvBytes, err := a.exportCandles.Export(context.Background(), instrument, timeframe, from, to)
	if err != nil {
		a.logger.Warn("export candles failed", "instrument", instrument, "error", err)
		_ = a.telegram.SendMessageToChat(context.Background(), chatID, "Error: "+err.Error())
		return
	}

	filename := fmt.Sprintf("candles_%s_%s.csv", instrument, timeframe)
	if err := a.telegram.SendDocument(context.Background(), chatID, filename, bytes.NewReader(csvBytes)); err != nil {
		a.logger.Error("SendDocument failed", "error", err)
		_ = a.telegram.SendMessageToChat(context.Background(), chatID, "Failed to send file: "+err.Error())
		return
	}
}

// parseCandlesArgs parses ["SBER", "1d", "30"] into instrument, timeframe, lastDays. Returns error if invalid.
func parseCandlesArgs(args []string) (instrument, timeframe string, lastDays int, err error) {
	if len(args) < 3 {
		return "", "", 0, fmt.Errorf("need instrument, timeframe and last_days (e.g. SBER 1d 30)")
	}
	instrument = strings.TrimSpace(args[0])
	timeframe = strings.TrimSpace(args[1])
	if instrument == "" || timeframe == "" {
		return "", "", 0, fmt.Errorf("instrument and timeframe must be non-empty")
	}
	var n int
	n, err = strconv.Atoi(strings.TrimSpace(args[2]))
	if err != nil || n <= 0 {
		return "", "", 0, fmt.Errorf("last_days must be a positive number")
	}
	return instrument, timeframe, n, nil
}
