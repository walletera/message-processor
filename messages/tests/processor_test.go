package tests

import (
    "sync"
    "testing"

    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "github.com/walletera/message-processor/events"
    "github.com/walletera/message-processor/messages"
)

func TestMessageProcessor_ProcessValidMessage(t *testing.T) {

    messagesCh := make(chan messages.Message)
    messageConsumerMock := &messages.MockConsumer{}
    messageConsumerMock.On("Consume").Return((<-chan messages.Message)(messagesCh), nil)

    rawPayload := []byte("raw message payload")
    message := messages.Message{Payload: rawPayload}

    event := FakeEvent{}

    eventsDeserializerMock := &events.MockDeserializer[FakeEventVisitor]{}
    eventsDeserializerMock.On("Deserialize", rawPayload).Return(event, nil)

    wg := sync.WaitGroup{}
    wg.Add(1)
    mockFakeEventVisitor := &MockFakeEventVisitor{}
    mockFakeEventVisitor.On("VisitFakeEvent", event).Return(nil).Run(func(args mock.Arguments) {
        wg.Done()
    })

    messageProcessor := messages.NewProcessor[FakeEventVisitor](
        messageConsumerMock,
        eventsDeserializerMock,
        mockFakeEventVisitor,
        func(processorError messages.ProcessorError) {},
    )

    messageProcessorStartError := messageProcessor.Start()
    require.NoError(t, messageProcessorStartError)

    messagesCh <- message

    messageConsumerMock.AssertExpectations(t)
    eventsDeserializerMock.AssertExpectations(t)
    mockFakeEventVisitor.AssertExpectations(t)
}
