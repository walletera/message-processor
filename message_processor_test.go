package main

import (
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "sync"
    "testing"
)

func TestMessageProcessor_ProcessValidMessage(t *testing.T) {

    messagesCh := make(chan Message)
    messageConsumerMock := &MockMessageConsumer{}
    messageConsumerMock.On("Consume").Return((<-chan Message)(messagesCh), nil)

    rawPayload := []byte("raw message payload")
    message := Message{Payload: rawPayload}

    event := WithdrawalCreated{
        Id:          "19cf3c4c-4a0d-417e-b0cc-83385b9487de",
        Amount:      100,
        Beneficiary: WithdrawalBeneficiary{},
    }

    eventsDeserializerMock := &MockEventsDeserializer[PaymentsEventsVisitor]{}
    eventsDeserializerMock.On("Deserialize", rawPayload).Return(event, nil)

    wg := sync.WaitGroup{}
    wg.Add(1)
    eventsHandlerMock := &MockPaymentsEventsVisitor{}
    eventsHandlerMock.On("VisitWithdrawalCreated", event).Return(nil).Run(func(args mock.Arguments) {
        wg.Done()
    })

    messageProcessor := NewMessageProcessor[PaymentsEventsVisitor](messageConsumerMock, eventsDeserializerMock, eventsHandlerMock)

    messageProcessorStartError := messageProcessor.Start()
    require.NoError(t, messageProcessorStartError)

    messagesCh <- message

    messageConsumerMock.AssertExpectations(t)
    eventsDeserializerMock.AssertExpectations(t)
    eventsHandlerMock.AssertExpectations(t)
}
