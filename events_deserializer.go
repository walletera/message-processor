package main

type EventsDeserializer interface {
    Deserialize(message Message) (any, error)
}
