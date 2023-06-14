package main

type EventsHandler interface {
    HandleEvent(event any) error
}
