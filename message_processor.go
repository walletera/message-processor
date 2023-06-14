package main

import (
    "fmt"
    "log"
)

type MessageProcessor struct {
    messageConsumer    MessageConsumer
    eventsDeserializer EventsDeserializer
    eventsHandler      EventsHandler
}

func NewMessageProcessor(
    messageConsumer MessageConsumer,
    eventsDeserializer EventsDeserializer,
    eventsHandler EventsHandler,
) *MessageProcessor {
    return &MessageProcessor{
        messageConsumer:    messageConsumer,
        eventsDeserializer: eventsDeserializer,
        eventsHandler:      eventsHandler,
    }
}

func (p *MessageProcessor) Start() error {
    msgCh, err := p.messageConsumer.Consume()
    if err != nil {
        return fmt.Errorf("failed consuming from message consumer: %w", err)
    }
    go p.processMessages(msgCh)
    return nil
}

func (p *MessageProcessor) processMessages(ch <-chan Message) {
    for message := range ch {
        go p.processMessage(message)
    }
}

func (p *MessageProcessor) processMessage(message Message) {
    event, err := p.eventsDeserializer.Deserialize(message)
    if err != nil {
        p.printErrorLog(message, err)
        return
    }
    if event != nil {
        err = p.eventsHandler.HandleEvent(event)
        if err != nil {
            p.printErrorLog(message, err)
        }
    }
}

func (p *MessageProcessor) printErrorLog(message Message, err error) {
    log.Printf("error processing message with payload: %s: %s", message.Payload, err.Error())
}
