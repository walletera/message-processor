package fake

import (
    "context"

    "github.com/walletera/message-processor/errors"
)

type EventVisitor interface {
    VisitFakeEvent(ctx context.Context, e Event) errors.ProcessingError
}
