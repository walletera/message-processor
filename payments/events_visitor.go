package payments

type EventsVisitor interface {
    VisitWithdrawalCreated(withdrawalCreated WithdrawalCreatedEvent) error
}
