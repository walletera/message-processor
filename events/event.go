package events

type EventData interface {
    ID() string
    Type() string
    CorrelationID() string
    DataContentType() string
    Serialize() ([]byte, error)
}

type Event[Visitor any] interface {
    EventData

    Accept(visitor Visitor) error
}
