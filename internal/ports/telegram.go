// Package ports defines interfaces (contracts) between app/domain and adapters.
// TelegramNotifier: send report text and optional PNG chart to Telegram. Implemented by adapters (e.g. Bot API).

package ports

import (
	"context"
	"io"
)

// PhotoItem is one image for SendPhotoAlbum (caption + reader; filename used for media group attach).
type PhotoItem struct {
	Caption  string    // caption for this photo
	Reader   io.Reader // image bytes (e.g. PNG)
	Filename string    // unique name for media group (e.g. "chart_0.png")
}

// TelegramNotifier sends report messages and optional images to a Telegram chat.
// Token and chat ID are configured in the adapter (e.g. from env). If not configured, adapter may no-op.
type TelegramNotifier interface {
	// SendMessage sends a text message to the configured chat.
	SendMessage(ctx context.Context, text string) error
	// SendMessageToChat sends a text message to a specific chat (e.g. for command replies or errors).
	SendMessageToChat(ctx context.Context, chatID int64, text string) error
	// SendPhoto sends a photo with optional caption. imageReader must contain PNG/JPEG bytes.
	SendPhoto(ctx context.Context, caption string, imageReader io.Reader) error
	// SendPhotoAlbum sends multiple photos as a media group (album). Max 10 per album; adapter may chunk.
	// If media group fails, adapter may fall back to sending each photo individually.
	SendPhotoAlbum(ctx context.Context, items []PhotoItem) error
	// SendDocument sends a file (e.g. CSV) to the given chat. Used for /candles command response.
	SendDocument(ctx context.Context, chatID int64, filename string, document io.Reader) error
	// ListenForMessages runs long polling and calls handler for each incoming text message.
	// Blocks until ctx is done. Used when handleCommands is enabled to process /candles etc.
	ListenForMessages(ctx context.Context, handler func(chatID int64, text string))
}
