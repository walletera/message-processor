package fake

import (
    "context"

    "github.com/walletera/message-processor/errors"
)

type EventHandler interface {
    HandleFakeEvent(ctx context.Context, e Event) errors.ProcessingError
}
