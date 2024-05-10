package main

import (
    "log"

    "github.com/walletera/message-processor/messages"
    "github.com/walletera/message-processor/payments"
)

func main() {

    paymentsEventsVisitor := NewPaymentsEventsVisitorImpl()

    processor, err := payments.NewRabbitMQProcessor(
        paymentsEventsVisitor,
        "message-processor-example-queue",
        func(err messages.ProcessorError) {
            log.Fatalf(err.Error())
        },
    )
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

func (p PaymentsEventsVisitorImpl) VisitWithdrawalCreated(withdrawalCreated payments.WithdrawalCreatedEvent) error {
    log.Printf("handling WithdrawalCreatedEvent event: %+v", withdrawalCreated)
    return nil
}
