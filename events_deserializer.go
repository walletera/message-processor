package main

type EventsDeserializer[Visitor any] interface {
    Deserialize(rawEvent []byte) (Event[Visitor], error)
}
