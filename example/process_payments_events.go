package main

import (
    "github.com/walletera/message-processor/pkg/events/payments"
    "github.com/walletera/message-processor/pkg/messages"
    "github.com/walletera/message-processor/pkg/rabbitmq"
    "log"
)

func main() {

    messageConsumer, err := rabbitmq.NewClient(rabbitmq.WithQueueName("message-consumer-queue"))
    if err != nil {
        log.Fatalf("failed initiating RabbitMQClient: %s", err.Error())
    }

    paymentsEventsDeserializer := payments.NewEventsDeserializer()
    paymentsEventsVisitor := NewPaymentsEventsVisitorImpl()

    processor := messages.NewProcessor[payments.EventsVisitor](messageConsumer, paymentsEventsDeserializer, paymentsEventsVisitor)

    err = processor.Start()
    if err != nil {
        log.Fatalf("failed to start message processor: %s", err.Error())
    }

    blockForeverCh := make(chan any)
    <-blockForeverCh
}

type PaymentsEventsVisitorImpl struct {
}

func NewPaymentsEventsVisitorImpl() *PaymentsEventsVisitorImpl {
    return &PaymentsEventsVisitorImpl{}
}

func (p PaymentsEventsVisitorImpl) VisitWithdrawalCreated(withdrawalCreated payments.WithdrawalCreated) error {
    log.Printf("handling WithdrawalCreated event: %+v", withdrawalCreated)
    return nil
}
