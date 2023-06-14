package main

import "log"

type PaymentsEventsVisitorImpl struct {
}

func NewPaymentsEventsVisitorImpl() *PaymentsEventsVisitorImpl {
    return &PaymentsEventsVisitorImpl{}
}

func (p PaymentsEventsVisitorImpl) VisitWithdrawalCreated(withdrawalCreated WithdrawalCreated) error {
    log.Printf("handling WithdrawalCreated event: %+v", withdrawalCreated)
    return nil
}
