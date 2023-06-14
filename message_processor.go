package main

import (
    "fmt"
    "log"
)

type MessageProcessor[Visitor any] struct {
    messageConsumer    MessageConsumer
    eventsDeserializer EventsDeserializer[Visitor]
    eventsVisitor      Visitor
}

func NewMessageProcessor[Visitor any](
    messageConsumer MessageConsumer,
    eventsDeserializer EventsDeserializer[Visitor],
    eventsVisitor Visitor,
) *MessageProcessor[Visitor] {

    return &MessageProcessor[Visitor]{
        messageConsumer:    messageConsumer,
        eventsDeserializer: eventsDeserializer,
        eventsVisitor:      eventsVisitor,
    }
}

func (p *MessageProcessor[Visitor]) Start() error {
    msgCh, err := p.messageConsumer.Consume()
    if err != nil {
        return fmt.Errorf("failed consuming from message consumer: %w", err)
    }
    go p.processMsgs(msgCh)
    return nil
}

func (p *MessageProcessor[Visitor]) processMsgs(ch <-chan Message) {
    for msg := range ch {
        go p.processMsg(msg)
    }
}

func (p *MessageProcessor[Visitor]) processMsg(message Message) {
    event, err := p.eventsDeserializer.Deserialize(message.Payload)
    if err != nil {
        p.printErrorLog(message, err)
        return
    }
    if event != nil {
        err = event.Accept(p.eventsVisitor)
        if err != nil {
            p.printErrorLog(message, err)
        }
    }
}

func (p *MessageProcessor[V]) printErrorLog(msg Message, err error) {
    log.Printf("error processing message with payload: %s: %s", msg.Payload, err.Error())
}
