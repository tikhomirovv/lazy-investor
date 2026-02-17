// Package ports defines interfaces (contracts) between app/domain and adapters.
// TelegramNotifier: send report text and optional PNG chart to Telegram. Implemented by adapters (e.g. Bot API).

package ports

import (
	"context"
	"io"
)

// TelegramNotifier sends report messages and optional images to a Telegram chat.
// Token and chat ID are configured in the adapter (e.g. from env). If not configured, adapter may no-op.
type TelegramNotifier interface {
	// SendMessage sends a text message to the configured chat.
	SendMessage(ctx context.Context, text string) error
	// SendPhoto sends a photo with optional caption. imageReader must contain PNG/JPEG bytes.
	SendPhoto(ctx context.Context, caption string, imageReader io.Reader) error
}
