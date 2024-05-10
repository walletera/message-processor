package messages

import (
    "fmt"

    "github.com/walletera/message-processor/events"
)

type Processor[Visitor any] struct {
    messageConsumer    Consumer
    eventsDeserializer events.Deserializer[Visitor]
    eventsVisitor      Visitor
    errorHandler       func(ProcessorError)
}

type ProcessorError struct {
    message Message
    error   error
}

func (p ProcessorError) Error() string {
    return fmt.Sprintf("error processing message with payload %s", p.message.Payload)
}

type ErrorHandler func(ProcessorError)

func NewProcessor[Visitor any](
    messageConsumer Consumer,
    eventsDeserializer events.Deserializer[Visitor],
    eventsVisitor Visitor,
    errorHandler ErrorHandler,
) *Processor[Visitor] {

    return &Processor[Visitor]{
        messageConsumer:    messageConsumer,
        eventsDeserializer: eventsDeserializer,
        eventsVisitor:      eventsVisitor,
        errorHandler:       errorHandler,
    }
}

func (p *Processor[Visitor]) Start() error {
    msgCh, err := p.messageConsumer.Consume()
    if err != nil {
        return fmt.Errorf("failed consuming from message consumer: %w", err)
    }
    go p.processMsgs(msgCh)
    return nil
}

func (p *Processor[Visitor]) processMsgs(ch <-chan Message) {
    for msg := range ch {
        go p.processMsg(msg)
    }
}

func (p *Processor[Visitor]) processMsg(message Message) {
    event, err := p.eventsDeserializer.Deserialize(message.Payload)
    if err != nil {
        p.errorHandler(ProcessorError{
            message: message,
            error:   err,
        })
        return
    }
    if event != nil {
        err = event.Accept(p.eventsVisitor)
        if err != nil {
            p.errorHandler(ProcessorError{
                message: message,
                error:   err,
            })
        }
    }
}
