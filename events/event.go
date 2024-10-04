package events

import (
	"context"

	"github.com/walletera/message-processor/errors"
)

type Event[Visitor any] interface {
	EventData

	Accept(ctx context.Context, visitor Visitor) errors.ProcessingError
}
