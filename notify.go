package notify

import (
	"context"

	"github.com/nmeilick/go-notify/message"
	"github.com/nmeilick/go-notify/registry"
	"github.com/nmeilick/go-notify/telegram"
)

func init() {
	registry.Register(&telegram.Telegram{})
}

func Notify(ctx context.Context, recipient string, m *message.Message) error {
	return registry.Notify(ctx, recipient, m)
}
