package main

type PaymentsEventsVisitor interface {
    VisitWithdrawalCreated(withdrawalCreated WithdrawalCreated) error
}
