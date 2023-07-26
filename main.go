package main

import (
    "log"
)

func main() {

    messageConsumer := NewRabbitMQMessageConsumer()
    paymentsEventsDeserializer := NewPaymentsEventsDeserializer()
    paymentsEventsVisitor := NewPaymentsEventsVisitorImpl()

    processor := NewMessageProcessor[PaymentsEventsVisitor](messageConsumer, paymentsEventsDeserializer, paymentsEventsVisitor)

    err := processor.Start()
    if err != nil {
        log.Fatalf("failed to start message processor: %s", err.Error())
    }

    blockForeverCh := make(chan any)
    <-blockForeverCh
}
