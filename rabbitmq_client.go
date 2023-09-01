package main

import (
    "context"
    "fmt"
    amqp "github.com/rabbitmq/amqp091-go"
)

const (
    DefaultPort      = 5672
    ManagementUIPort = 15672

    ExchangeTypeDirect = "direct"
    ExchangeTypeTopic  = "topic"
    ExchangeTypeFanout = "fanout"
)

type RabbitMQClient struct {
    conn        *amqp.Connection
    connChannel *amqp.Channel
    queue       amqp.Queue

    port                uint
    useDefaultExchange  bool
    exchangeName        string
    exchangeType        string
    queueName           string
    consumerRoutingKeys []string
}

type ConsumerOpt func(consumer *RabbitMQClient)

func NewRabbitMQClient(opts ...ConsumerOpt) (*RabbitMQClient, error) {
    consumer := &RabbitMQClient{}
    err := applyOptionsOrDefault(consumer, opts)
    if err != nil {
        return nil, err
    }
    err = consumer.init()
    if err != nil {
        return nil, err
    }
    return consumer, nil
}

func (r *RabbitMQClient) Consume() (<-chan Message, error) {
    if r.connChannel == nil {
        return nil, fmt.Errorf("RabbitMQClient was not properly initialized")
    }

    if !r.useDefaultExchange {
        if len(r.consumerRoutingKeys) == 0 {
            return nil, fmt.Errorf("missing routing key")
        }
        for _, key := range r.consumerRoutingKeys {
            err := r.connChannel.QueueBind(
                r.queue.Name,   // queue name
                key,            // routing key
                r.exchangeName, // exchange
                false,
                nil)
            if err != nil {
                fmt.Errorf("failed to bind queue %s with exchange %s using routing key %s", r.queue.Name, r.exchangeName, key)
            }
        }
    }

    msgs, err := r.connChannel.Consume(
        r.queue.Name, // queue
        "",           // consumer
        true,         // auto-ack
        false,        // exclusive
        false,        // no-local
        false,        // no-wait
        nil,          // args
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

func (r *RabbitMQClient) Publish(ctx context.Context, message []byte, routingKey string) error {

    err := r.connChannel.PublishWithContext(ctx,
        r.exchangeName, // exchange
        routingKey,     // routing key
        false,          // mandatory
        false,          // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        message,
        })
    if err != nil {
        return err
    }

    return nil
}

func (r *RabbitMQClient) Close() error {
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

func (r *RabbitMQClient) QueueName() string {
    return r.queue.Name
}

func (r *RabbitMQClient) init() error {
    conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@localhost:%d/", r.port))
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    r.conn = conn

    ch, err := conn.Channel()
    if err != nil {
        return fmt.Errorf("failed to open a channel: %w", err)
    }

    r.connChannel = ch

    if !r.useDefaultExchange {
        err = ch.ExchangeDeclare(
            r.exchangeName, // name
            r.exchangeType, // type
            true,           // durable
            false,          // auto-deleted
            false,          // internal
            false,          // no-wait
            nil,            // arguments
        )
        if err != nil {
            return fmt.Errorf("failed to declare exchange: %w", err)
        }
    }

    q, err := ch.QueueDeclare(
        r.queueName, // name
        false,       // durable
        false,       // delete when unused
        false,       // exclusive
        false,       // no-wait
        nil,         // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare a queue: %w", err)
    }

    r.queue = q

    return nil
}

func applyOptionsOrDefault(consumer *RabbitMQClient, opts []ConsumerOpt) error {
    consumer.port = DefaultPort
    consumer.useDefaultExchange = true
    for _, opt := range opts {
        opt(consumer)
    }
    if !consumer.useDefaultExchange {
        if consumer.exchangeType == "" {
            return fmt.Errorf("if useDefaultExchange is false exchange type can't be empty")
        }
    }
    return nil
}

func WithPort(port uint) func(c *RabbitMQClient) {
    return func(c *RabbitMQClient) {
        c.port = port
    }
}

func WithQueueName(queueName string) func(c *RabbitMQClient) {
    return func(c *RabbitMQClient) {
        c.queueName = queueName
    }
}

func WithExchangeName(exchangeName string) func(c *RabbitMQClient) {
    return func(c *RabbitMQClient) {
        c.useDefaultExchange = false
        c.exchangeName = exchangeName
    }
}

func WithExchangeType(exchangeType string) func(c *RabbitMQClient) {
    return func(c *RabbitMQClient) {
        c.exchangeType = exchangeType
    }
}

func WithConsumerRoutingKeys(routingKeys ...string) func(c *RabbitMQClient) {
    return func(c *RabbitMQClient) {
        c.consumerRoutingKeys = routingKeys
    }
}
