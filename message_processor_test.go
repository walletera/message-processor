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

    message := Message{
        Payload: []byte{},
    }

    event := WithdrawalCreated{
        WithdrawalId:       "19cf3c4c-4a0d-417e-b0cc-83385b9487de",
        Amount:             100,
        SourceAccount:      "source account details",
        DestinationAccount: "destination account details",
    }

    eventsDeserializerMock := &MockEventsDeserializer{}
    eventsDeserializerMock.On("Deserialize", message).Return(event, nil)

    wg := sync.WaitGroup{}
    wg.Add(1)
    eventsHandlerMock := &MockEventsHandler{}
    eventsHandlerMock.On("HandleEvent", event).Return(nil).Run(func(args mock.Arguments) {
        wg.Done()
    })

    messageProcessor := NewMessageProcessor(messageConsumerMock, eventsDeserializerMock, eventsHandlerMock)

    messageProcessorStartError := messageProcessor.Start()
    require.NoError(t, messageProcessorStartError)

    messagesCh <- message

    messageConsumerMock.AssertExpectations(t)
    eventsDeserializerMock.AssertExpectations(t)
    eventsHandlerMock.AssertExpectations(t)
}
