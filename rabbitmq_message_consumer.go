package main

import (
    "fmt"
    amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQMessageConsumer struct {
    conn        *amqp.Connection
    connChannel *amqp.Channel
}

func NewRabbitMQMessageConsumer() *RabbitMQMessageConsumer {
    return &RabbitMQMessageConsumer{}
}

func (r *RabbitMQMessageConsumer) Consume() (<-chan Message, error) {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    r.conn = conn

    ch, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open a channel: %w", err)
    }

    r.connChannel = ch

    q, err := ch.QueueDeclare(
        "message-consumer-queue", // name
        false,                    // durable
        false,                    // delete when unused
        false,                    // exclusive
        false,                    // no-wait
        nil,                      // arguments
    )
    if err != nil {
        return nil, fmt.Errorf("failed to declare a queue: %w", err)
    }

    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    if err != nil {
        return nil, fmt.Errorf("failed to register a consumer: %w", err)
    }

    messagesCh := make(chan Message)
    go func() {
        defer close(messagesCh)
        for msg := range msgs {
            messagesCh <- Message{
                Payload: msg.Body,
            }
        }
    }()

    return messagesCh, nil
}

func (r *RabbitMQMessageConsumer) Close() error {
    err := r.connChannel.Close()
    if err != nil {
        return fmt.Errorf("failed to close rabbitmq connection channel: %w", err)
    }
    err = r.conn.Close()
    if err != nil {
        return fmt.Errorf("failed to close rabbitmq connection: %w", err)
    }
    return nil
}
