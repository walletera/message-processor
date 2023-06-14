package main

import (
    "log"
)

func main() {

    messageConsumer := NewRabbitMQMessageConsumer()
    paymentsEventsDeserializer := NewPaymentsEventsDeserializerImpl()
    paymentsEventsHandler := NewPaymentsEventsHandler()

    processor := NewMessageProcessor(messageConsumer, paymentsEventsDeserializer, paymentsEventsHandler)

    err := processor.Start()
    if err != nil {
        log.Fatalf("failed to start message processor: %s", err.Error())
    }

    blockForeverCh := make(chan any)
    <-blockForeverCh
}
