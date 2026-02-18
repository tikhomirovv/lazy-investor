// Package telegram implements ports.TelegramNotifier using go-telegram-bot-api/v5.
// Token and ChatID come from env; if either is missing, methods no-op and return nil.
package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tikhomirovv/lazy-investor/internal/ports"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

const requestTimeout = 15 * time.Second

// Ensure Service implements ports.TelegramNotifier.
var _ ports.TelegramNotifier = (*Service)(nil)

// Config holds Telegram connection settings (from env). Empty token or chatID means no-op.
type Config struct {
	Token  string // TELEGRAM_BOT_TOKEN
	ChatID string // TELEGRAM_CHAT_ID (numeric: your user ID for private chat, or group ID)
}

// Service sends messages and photos via go-telegram-bot-api.
type Service struct {
	config Config
	logger logging.Logger
	bot    *tgbotapi.BotAPI
}

// NewService creates the Telegram adapter. Pass token and chatID from env; empty = no-op mode.
// Bot API client is created only when both token and chatID are set; invalid chatID (non-numeric) is ignored (no-op).
func NewService(config Config, logger logging.Logger) *Service {
	if config.Token == "" || config.ChatID == "" {
		return &Service{config: config, logger: logger, bot: nil}
	}
	client := &http.Client{Timeout: requestTimeout}
	bot, err := tgbotapi.NewBotAPIWithClient(config.Token, tgbotapi.APIEndpoint, client)
	if err != nil {
		logger.Warn("Telegram bot init failed (will no-op)", "error", err)
		return &Service{config: config, logger: logger, bot: nil}
	}
	return &Service{config: config, logger: logger, bot: bot}
}

func (s *Service) chatIDInt64() (int64, bool) {
	if s.config.ChatID == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(s.config.ChatID, 10, 64)
	if err != nil {
		s.logger.Warn("TELEGRAM_CHAT_ID must be numeric (Bot API does not accept @username). Use your user ID from @userinfobot or @getmyid_bot", "chat_id", s.config.ChatID)
		return 0, false
	}
	return id, true
}

// SendMessage sends a text message. No-op if bot or chatID not configured.
func (s *Service) SendMessage(ctx context.Context, text string) error {
	if s.bot == nil {
		s.logger.Debug("Telegram not configured (missing token or chat_id), skipping SendMessage")
		return nil
	}
	chatID, ok := s.chatIDInt64()
	if !ok {
		return nil
	}
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := s.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("telegram SendMessage: %w", err)
	}
	return nil
}

// SendPhoto sends a photo with optional caption. imageReader must contain image bytes (e.g. PNG).
// No-op if bot or chatID not configured.
func (s *Service) SendPhoto(ctx context.Context, caption string, imageReader io.Reader) error {
	if s.bot == nil {
		s.logger.Debug("Telegram not configured (missing token or chat_id), skipping SendPhoto")
		return nil
	}
	chatID, ok := s.chatIDInt64()
	if !ok {
		return nil
	}
	file := tgbotapi.FileReader{Name: "chart.png", Reader: imageReader}
	photo := tgbotapi.NewPhoto(chatID, file)
	photo.Caption = caption
	_, err := s.bot.Send(photo)
	if err != nil {
		return fmt.Errorf("telegram SendPhoto: %w", err)
	}
	return nil
}
