package fake

import (
    "context"

    "github.com/walletera/message-processor/errors"
    "github.com/walletera/message-processor/events"
)

var _ events.Event[EventVisitor] = Event{}

type Event struct {
}

func (f Event) ID() string {
    //TODO implement me
    panic("implement me")
}

func (f Event) Type() string {
    //TODO implement me
    panic("implement me")
}

func (f Event) CorrelationID() string {
    //TODO implement me
    panic("implement me")
}

func (f Event) DataContentType() string {
    //TODO implement me
    panic("implement me")
}

func (f Event) Serialize() ([]byte, error) {
    //TODO implement me
    panic("implement me")
}

func (f Event) Accept(ctx context.Context, visitor EventVisitor) errors.ProcessingError {
    return visitor.VisitFakeEvent(ctx, f)
}
