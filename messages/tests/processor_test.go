package tests

import (
    "fmt"
    "sync"
    "testing"

    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "github.com/walletera/message-processor/errors"
    "github.com/walletera/message-processor/fake"
    "github.com/walletera/message-processor/messages"
    eventsmock "github.com/walletera/message-processor/mocks/github.com/walletera/message-processor/events"
    fakemock "github.com/walletera/message-processor/mocks/github.com/walletera/message-processor/fake"
    messagesmock "github.com/walletera/message-processor/mocks/github.com/walletera/message-processor/messages"
)

func TestMessageProcessor_ProcessValidMessage(t *testing.T) {

    var testErrs = map[string]errors.ProcessingError{
        "noErr":                     nil,
        "unprocessableMessageError": errors.NewUnprocessableMessageError("boom"),
        "internalError":             errors.NewInternalError("boom"),
        "timeoutError":              errors.NewTimeoutError("boom"),
    }

    for errName, err := range testErrs {
        var acknowledgerMockExpectationSetter func(mock *messagesmock.MockAcknowledger)
        var testName string
        if errName == "noErr" {
            acknowledgerMockExpectationSetter = func(mock *messagesmock.MockAcknowledger) {
                mock.On("Ack").Return(nil)
            }
            testName = "if no processing error then the message is Ack"
        } else {
            acknowledgerMockExpectationSetter = func(mock *messagesmock.MockAcknowledger) {
                mock.On("Nack", messages.NackOpts{
                    Requeue:      err.IsRetryable(),
                    ErrorCode:    err.Code(),
                    ErrorMessage: err.Message(),
                }).Return(nil)
            }
            testName = fmt.Sprintf("if %s occurred then message is Nack", errName)
        }

        t.Run(testName, func(t *testing.T) {
            execTest(t, acknowledgerMockExpectationSetter, err)
        })
    }
}

func execTest(t *testing.T, acknowledgerMockExpectationSetter func(acknowledgerMock *messagesmock.MockAcknowledger), processingErr errors.ProcessingError) {
    messagesCh := make(chan messages.Message)
    messageConsumerMock := &messagesmock.MockConsumer{}
    messageConsumerMock.On("Consume").Return((<-chan messages.Message)(messagesCh), nil)

    acknowledgerMock := &messagesmock.MockAcknowledger{}
    acknowledgerMockExpectationSetter(acknowledgerMock)

    rawPayload := []byte("raw message payload")
    message := messages.NewMessage(rawPayload, acknowledgerMock)

    event := fake.Event{}

    eventsDeserializerMock := &eventsmock.MockDeserializer[fake.EventVisitor]{}
    eventsDeserializerMock.On("Deserialize", rawPayload).Return(event, nil)

    wg := sync.WaitGroup{}
    wg.Add(1)
    mockFakeEventVisitor := &fakemock.MockEventVisitor{}
    mockFakeEventVisitor.
        On("VisitFakeEvent", mock.Anything, event).
        Return(processingErr).
        Run(func(args mock.Arguments) {
            wg.Done()
        })

    messageProcessor := messages.NewProcessor[fake.EventVisitor](
        messageConsumerMock,
        eventsDeserializerMock,
        mockFakeEventVisitor,
    )

    messageProcessorStartError := messageProcessor.Start()
    require.NoError(t, messageProcessorStartError)

    messagesCh <- message

    wg.Wait()

    acknowledgerMock.AssertExpectations(t)
    messageConsumerMock.AssertExpectations(t)
    eventsDeserializerMock.AssertExpectations(t)
    mockFakeEventVisitor.AssertExpectations(t)
}
