package payments

import (
    "fmt"

    "github.com/walletera/message-processor/messages"
    "github.com/walletera/message-processor/rabbitmq"
)

const (
    RabbitMQExchangeName = "payments"
    RabbitMQExchangeType = "topic"
    RabbitMQRoutingKey   = "payments.events"
)

func NewRabbitMQProcessor(
    eventsVisitor EventsVisitor,
    queueName string,
    errorHandler messages.ErrorHandler,
    rabbitmqOpts ...rabbitmq.ConsumerOpt,
) (*messages.Processor[EventsVisitor], error) {
    defaultOpts := []rabbitmq.ConsumerOpt{
        rabbitmq.WithExchangeName(RabbitMQExchangeName),
        rabbitmq.WithExchangeType(RabbitMQExchangeType),
        rabbitmq.WithConsumerRoutingKeys(RabbitMQRoutingKey),
        rabbitmq.WithQueueName(queueName),
    }
    opts := append(defaultOpts, rabbitmqOpts...)
    rabbitMQClient, err := rabbitmq.NewClient(opts...)
    if err != nil {
        return nil, fmt.Errorf("creating rabbitmq client: %w", err)
    }
    return NewProcessor(rabbitMQClient, eventsVisitor, errorHandler), nil
}
