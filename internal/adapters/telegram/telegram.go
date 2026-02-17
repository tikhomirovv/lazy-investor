// Package telegram implements ports.TelegramNotifier via Telegram Bot API (HTTP).
// Token and ChatID come from env; if either is missing, methods no-op and return nil.
package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/tikhomirovv/lazy-investor/internal/ports"
	"github.com/tikhomirovv/lazy-investor/pkg/logging"
)

const (
	// apiBase is the Telegram Bot API base URL (no trailing slash).
	apiBase     = "https://api.telegram.org"
	httpTimeout = 15 * time.Second
)

// Ensure Service implements ports.TelegramNotifier.
var _ ports.TelegramNotifier = (*Service)(nil)

// Config holds Telegram connection settings (from env). Empty token or chatID means no-op.
type Config struct {
	Token  string // TELEGRAM_BOT_TOKEN
	ChatID string // TELEGRAM_CHAT_ID
}

// Service sends messages and photos via Telegram Bot API.
type Service struct {
	config Config
	logger logging.Logger
	client *http.Client
}

// NewService creates the Telegram adapter. Pass token and chatID from env; empty = no-op mode.
func NewService(config Config, logger logging.Logger) *Service {
	return &Service{
		config: config,
		logger: logger,
		client: &http.Client{Timeout: httpTimeout},
	}
}

// SendMessage sends a text message. No-op if token or chatID not configured.
func (s *Service) SendMessage(ctx context.Context, text string) error {
	if s.config.Token == "" || s.config.ChatID == "" {
		s.logger.Debug("Telegram not configured (missing token or chat_id), skipping SendMessage")
		return nil
	}
	u := fmt.Sprintf("%s/bot%s/sendMessage", apiBase, s.config.Token)
	body := url.Values{}
	body.Set("chat_id", s.config.ChatID)
	body.Set("text", text)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBufferString(body.Encode()))
	if err != nil {
		return fmt.Errorf("telegram sendMessage request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram sendMessage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bs, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram sendMessage: status %d, body: %s", resp.StatusCode, string(bs))
	}
	var out struct {
		OK bool `json:"ok"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("telegram sendMessage decode: %w", err)
	}
	if !out.OK {
		return fmt.Errorf("telegram sendMessage: api returned ok=false")
	}
	return nil
}

// SendPhoto sends a photo with optional caption. imageReader must contain image bytes (e.g. PNG).
// No-op if token or chatID not configured.
func (s *Service) SendPhoto(ctx context.Context, caption string, imageReader io.Reader) error {
	if s.config.Token == "" || s.config.ChatID == "" {
		s.logger.Debug("Telegram not configured (missing token or chat_id), skipping SendPhoto")
		return nil
	}
	u := fmt.Sprintf("%s/bot%s/sendPhoto", apiBase, s.config.Token)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("chat_id", s.config.ChatID)
	if caption != "" {
		_ = w.WriteField("caption", caption)
	}
	part, err := w.CreateFormFile("photo", "chart.png")
	if err != nil {
		return fmt.Errorf("telegram sendPhoto form: %w", err)
	}
	if _, err := io.Copy(part, imageReader); err != nil {
		return fmt.Errorf("telegram sendPhoto copy: %w", err)
	}
	contentType := w.FormDataContentType()
	if err := w.Close(); err != nil {
		return fmt.Errorf("telegram sendPhoto close: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, &buf)
	if err != nil {
		return fmt.Errorf("telegram sendPhoto request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram sendPhoto: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bs, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram sendPhoto: status %d, body: %s", resp.StatusCode, string(bs))
	}
	var out struct {
		OK bool `json:"ok"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return fmt.Errorf("telegram sendPhoto decode: %w", err)
	}
	if !out.OK {
		return fmt.Errorf("telegram sendPhoto: api returned ok=false")
	}
	return nil
}
