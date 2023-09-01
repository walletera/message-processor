package main

import (
    "log"
)

func main() {

    messageConsumer, err := NewRabbitMQClient(WithQueueName("message-consumer-queue"))
    if err != nil {
        log.Fatalf("failed initiating RabbitMQClient: %s", err.Error())
    }

    paymentsEventsDeserializer := NewPaymentsEventsDeserializer()
    paymentsEventsVisitor := NewPaymentsEventsVisitorImpl()

    processor := NewMessageProcessor[PaymentsEventsVisitor](messageConsumer, paymentsEventsDeserializer, paymentsEventsVisitor)

    err = processor.Start()
    if err != nil {
        log.Fatalf("failed to start message processor: %s", err.Error())
    }

    blockForeverCh := make(chan any)
    <-blockForeverCh
}
