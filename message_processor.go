package main

import (
    "fmt"
    "log"
)

type MessageProcessor struct {
    messageConsumer *RabbitMQMessageConsumer
}

func NewMessageProcessor() *MessageProcessor {
    return &MessageProcessor{
        messageConsumer: NewRabbitMQMessageConsumer(),
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

func (p *MessageProcessor) processMessage(msg Message) {
    log.Printf("processing message with payload %s", msg.Payload)
}
