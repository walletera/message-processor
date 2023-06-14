package main

type MessageConsumer interface {
    Consume() (<-chan Message, error)
    Close() error
}
