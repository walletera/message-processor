package events

type EventData interface {
    ID() string
    Type() string
    DataContentType() string
    Serialize() ([]byte, error)
}

type Event[Visitor any] interface {
    EventData

    Accept(visitor Visitor) error
}
