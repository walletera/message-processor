package main

import (
    "github.com/walletera/message-processor/pkg/events/payments"
    "log"
)

func main() {

    paymentsEventsVisitor := NewPaymentsEventsVisitorImpl()

    processor, err := payments.NewRabbitMQProcessor(paymentsEventsVisitor, "message-processor-example-queue")
    if err != nil {
        log.Fatalf("error creating new RabbitMQProcessor: %s", err.Error())
    }

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
