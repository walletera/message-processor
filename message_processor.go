package main

import (
    "fmt"
    "log"
)

type MessageProcessor struct {
    messageConsumer            *RabbitMQMessageConsumer
    paymentsEventsDeserializer *PaymentsEventsDeserializer
    paymentsEventsHandler      *PaymentsEventsHandler
}

func NewMessageProcessor() *MessageProcessor {
    return &MessageProcessor{
        messageConsumer:            NewRabbitMQMessageConsumer(),
        paymentsEventsDeserializer: NewPaymentsEventsDeserializer(),
        paymentsEventsHandler:      NewPaymentsEventsHandler(),
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
    for msg := range ch {
        go p.processMessage(msg)
    }
}

func (p *MessageProcessor) processMessage(message Message) {
    event, err := p.paymentsEventsDeserializer.Deserialize(message)
    if err != nil {
        p.pringErrorLog(message, err)
    }
    err = p.paymentsEventsHandler.HandleEvent(event)
    if err != nil {
        p.pringErrorLog(message, err)
    }
}

func (p *MessageProcessor) pringErrorLog(message Message, err error) {
    log.Printf("error processing message %s: %s", message.Payload, err.Error())
}
