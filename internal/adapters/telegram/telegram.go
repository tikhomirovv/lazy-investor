// Package telegram implements ports.TelegramNotifier using go-telegram-bot-api/v5.
// Token and ChatID come from env; if either is missing, methods no-op and return nil.
package telegram

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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

// MaxMediaGroupSize is Telegram's limit for one sendMediaGroup request.
const MaxMediaGroupSize = 10

// SendPhotoAlbum sends multiple photos as media group(s). Chunks by MaxMediaGroupSize (10).
// Reads each item.Reader into memory so that on media group failure we can fall back to SendPhoto per image.
func (s *Service) SendPhotoAlbum(ctx context.Context, items []ports.PhotoItem) error {
	if s.bot == nil {
		s.logger.Debug("Telegram not configured (missing token or chat_id), skipping SendPhotoAlbum")
		return nil
	}
	if len(items) == 0 {
		return nil
	}
	chatID, ok := s.chatIDInt64()
	if !ok {
		return nil
	}
	// Read all into memory so fallback can resend (readers are one-shot).
	var all []photoData
	for i := range items {
		data, err := io.ReadAll(items[i].Reader)
		if err != nil {
			return fmt.Errorf("telegram SendPhotoAlbum read item %d: %w", i, err)
		}
		name := items[i].Filename
		if name == "" {
			name = fmt.Sprintf("chart_%d.png", i)
		}
		all = append(all, photoData{caption: items[i].Caption, filename: name, data: data})
	}
	for start := 0; start < len(all); start += MaxMediaGroupSize {
		end := start + MaxMediaGroupSize
		if end > len(all) {
			end = len(all)
		}
		chunk := all[start:end]
		if err := s.sendMediaGroupChunk(chatID, chunk); err != nil {
			// Library fails to unmarshal media group response (API returns array, lib expects single Message).
			// When that happens, the album was actually sent; do not fall back or we send duplicates.
			if isMediaGroupResponseUnmarshalError(err) {
				s.logger.Debug("SendPhotoAlbum: media group sent (library unmarshal quirk ignored)")
				continue
			}
			s.logger.Warn("SendPhotoAlbum: media group failed, falling back to individual SendPhoto", "error", err)
			for i := range chunk {
				if e := s.SendPhoto(ctx, chunk[i].caption, bytes.NewReader(chunk[i].data)); e != nil {
					return fmt.Errorf("telegram SendPhoto fallback: %w", e)
				}
			}
		}
	}
	return nil
}

// photoData holds in-memory image bytes for one album item (used for media group and fallback).
type photoData struct {
	caption  string
	filename string
	data     []byte
}

// isMediaGroupResponseUnmarshalError returns true when the error is the known go-telegram-bot-api
// quirk: Send(mediaGroup) succeeds on the wire but the library fails to unmarshal the response
// (API returns []Message, library expects Message). In that case the album was delivered.
func isMediaGroupResponseUnmarshalError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "cannot unmarshal array into Go value of type") && strings.Contains(s, "Message")
}

// sendMediaGroupChunk sends one album of up to 10 photos. chunk contains already-read data.
func (s *Service) sendMediaGroupChunk(chatID int64, chunk []photoData) error {
	media := make([]interface{}, 0, len(chunk))
	for i := range chunk {
		inputMedia := tgbotapi.NewInputMediaPhoto(tgbotapi.FileBytes{Name: chunk[i].filename, Bytes: chunk[i].data})
		inputMedia.Caption = chunk[i].caption
		media = append(media, inputMedia)
	}
	cfg := tgbotapi.NewMediaGroup(chatID, media)
	_, err := s.bot.Send(cfg)
	if err != nil {
		return fmt.Errorf("telegram SendMediaGroup: %w", err)
	}
	return nil
}
