package main

import (
    "log"
)

func main() {

    messageConsumer := NewRabbitMQMessageConsumer()
    paymentsEventsDeserializer := NewPaymentsEventsDeserializer()
    paymentsEventsHandler := NewPaymentsEventsVisitorImpl()

    processor := NewMessageProcessor[PaymentsEventsVisitor](messageConsumer, paymentsEventsDeserializer, paymentsEventsHandler)

    err := processor.Start()
    if err != nil {
        log.Fatalf("failed to start message processor: %s", err.Error())
    }

    blockForeverCh := make(chan any)
    <-blockForeverCh
}
