package messages

import (
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "github.com/walletera/message-processor/pkg/events"
    "github.com/walletera/message-processor/pkg/events/payments"
    "sync"
    "testing"
)

func TestMessageProcessor_ProcessValidMessage(t *testing.T) {

    messagesCh := make(chan Message)
    messageConsumerMock := &MockConsumer{}
    messageConsumerMock.On("Consume").Return((<-chan Message)(messagesCh), nil)

    rawPayload := []byte("raw message payload")
    message := Message{Payload: rawPayload}

    event := payments.WithdrawalCreated{
        Id:          "19cf3c4c-4a0d-417e-b0cc-83385b9487de",
        Amount:      100,
        Beneficiary: payments.WithdrawalBeneficiary{},
    }

    eventsDeserializerMock := &events.MockDeserializer[payments.EventsVisitor]{}
    eventsDeserializerMock.On("Deserialize", rawPayload).Return(event, nil)

    wg := sync.WaitGroup{}
    wg.Add(1)
    eventsHandlerMock := &payments.MockEventsVisitor{}
    eventsHandlerMock.On("VisitWithdrawalCreated", event).Return(nil).Run(func(args mock.Arguments) {
        wg.Done()
    })

    messageProcessor := NewProcessor[payments.EventsVisitor](messageConsumerMock, eventsDeserializerMock, eventsHandlerMock)

    messageProcessorStartError := messageProcessor.Start()
    require.NoError(t, messageProcessorStartError)

    messagesCh <- message

    messageConsumerMock.AssertExpectations(t)
    eventsDeserializerMock.AssertExpectations(t)
    eventsHandlerMock.AssertExpectations(t)
}
