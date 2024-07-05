package main

import (
    "context"
    "log"

    "github.com/walletera/message-processor/errors"
    "github.com/walletera/message-processor/messages"
    "github.com/walletera/message-processor/payments"
)

func main() {

    paymentsEventsVisitor := NewPaymentsEventsVisitorImpl()

    processor, err := payments.NewRabbitMQProcessor(
        paymentsEventsVisitor,
        "message-processor-example-queue",
        payments.RabbitMQProcessorOpt{
            ProcessorOpt: messages.WithErrorCallback(
                func(processingError errors.ProcessingError) {
                    log.Printf("error processing message: %s", processingError.Error())
                },
            ),
        },
    )
    if err != nil {
        log.Fatalf("error creating new RabbitMQProcessor: %s", err.Error())
    }

    err = processor.Start(context.Background())
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

func (p PaymentsEventsVisitorImpl) VisitWithdrawalCreated(_ context.Context, withdrawalCreated payments.WithdrawalCreatedEvent) errors.ProcessingError {
    log.Printf("handling WithdrawalCreatedEvent event: %+v", withdrawalCreated)
    return nil
}
