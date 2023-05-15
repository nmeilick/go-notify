// Registry manages notifiers and serves as the central entry point to send notifications.
package registry

import (
	"context"
	"sync"

	"github.com/nmeilick/go-notify/message"
)

// Notifier is an interface that represents a notification service such as email, SMS, push notifications, etc.
type Notifier interface {
	Notify(context.Context, string, *message.Message) error
}

var mu sync.RWMutex
var notifiers []Notifier

// Register adds one or more Notifier instances to the registry.
func Register(nn ...Notifier) {
	mu.Lock()
	defer mu.Unlock()
	notifiers = append(notifiers, nn...)
}

// Notify sends a notification to the specified recipient using one of the registered Notifier instances.
func Notify(ctx context.Context, recipient string, m *message.Message) error {
	// Create a local copy of notifiers
	mu.RLock()
	notifiers := append([]Notifier(nil), notifiers...)
	mu.RUnlock()

	for _, notifier := range notifiers {
		if err := notifier.Notify(ctx, recipient, m); err != nil {
			return err
		}
	}
	return nil
}
