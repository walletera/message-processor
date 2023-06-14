package main

import (
    "fmt"
    "log"
)

type PaymentsEventsHandler struct {
}

func NewPaymentsEventsHandler() *PaymentsEventsHandler {
    return &PaymentsEventsHandler{}
}

func (h *PaymentsEventsHandler) HandleEvent(event any) error {
    switch event.(type) {
    case WithdrawalCreated:
        log.Printf("handling WithdrawalCreated event: %+v", event)
    default:
        return fmt.Errorf("unexpected event type %t", event)
    }
    return nil
}
