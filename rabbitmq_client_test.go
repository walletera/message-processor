//go:build rabbitmq_client_test

package main

import (
    "context"
    "fmt"
    "github.com/stretchr/testify/assert"
    "math/rand"
    "os"
    "testing"
    "time"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func TestRabbitMQClientWithDefaults(t *testing.T) {
    ctx := context.Background()

    stopRabbitMQ := runRabbitMQ(ctx)
    defer stopRabbitMQ(ctx)

    consumer, err := NewRabbitMQClient()
    if err != nil {
        t.Error(err.Error())
        return
    }
    defer consumer.Close()

    messagesCh, err := consumer.Consume()
    if err != nil {
        t.Errorf("error executing consumer.Consumer: %s", err.Error())
        return
    }

    randomNum := rand.New(rand.NewSource(time.Now().Unix())).Int63()
    message := fmt.Sprintf("a message with a random number %d", randomNum)

    publisher, err := NewRabbitMQClient()
    if err != nil {
        t.Error(err.Error())
        return
    }
    defer publisher.Close()

    err = publisher.Publish(ctx, []byte(message), consumer.QueueName())
    if err != nil {
        t.Errorf("error publishing message: %s", err.Error())
        return
    }

    waitForMessageWithTimeout(t, message, messagesCh, 5*time.Second)
}

func TestRabbitMQClientWithOptions(t *testing.T) {
    ctx := context.Background()
    stopRabbitMQ := runRabbitMQ(ctx)

    defer stopRabbitMQ(ctx)

    const testRoutingKey = "test.routing.key"

    consumer, err := NewRabbitMQClient(
        WithQueueName("test-queue"),
        WithExchangeName("test-exchange"),
        WithExchangeType(ExchangeTypeTopic),
        WithConsumerRoutingKeys(testRoutingKey),
    )
    if err != nil {
        t.Errorf("error creating RabbitMQClient: %s", err.Error())
        return
    }
    defer consumer.Close()

    messagesCh, err := consumer.Consume()
    if err != nil {
        t.Errorf("error executing consumer.Consumer: %s", err.Error())
        return
    }

    randomNum := rand.New(rand.NewSource(time.Now().Unix())).Int63()
    message := fmt.Sprintf("a message with a random number %d", randomNum)

    publisher, err := NewRabbitMQClient(
        WithQueueName("test-queue"),
        WithExchangeName("test-exchange"),
        WithExchangeType(ExchangeTypeTopic),
        WithConsumerRoutingKeys(testRoutingKey),
    )
    if err != nil {
        t.Error(err.Error())
        return
    }
    defer publisher.Close()

    err = publisher.Publish(ctx, []byte(message), testRoutingKey)
    if err != nil {
        t.Errorf("error publishing message: %s", err.Error())
        return
    }

    waitForMessageWithTimeout(t, message, messagesCh, 5*time.Second)
}

func TestRabbitMQClientWithOptionsButEmptyExchangeType(t *testing.T) {
    ctx := context.Background()
    stopRabbitMQ := runRabbitMQ(ctx)

    defer stopRabbitMQ(ctx)

    const testRoutingKey = "test.routing.key"

    _, err := NewRabbitMQClient(
        WithQueueName("test-queue"),
        WithExchangeName("test-exchange"),
        WithConsumerRoutingKeys(testRoutingKey),
    )
    assert.Error(t, err)
}

func TestRabbitMQClientWithLoad(t *testing.T) {
    ctx := context.Background()
    stopRabbitMQ := runRabbitMQ(ctx)

    defer stopRabbitMQ(ctx)

    const testRoutingKey = "test.routing.key"

    consumer, err := NewRabbitMQClient(
        WithQueueName("test-queue"),
        WithExchangeName("test-exchange"),
        WithExchangeType(ExchangeTypeTopic),
        WithConsumerRoutingKeys(testRoutingKey),
    )
    if err != nil {
        t.Errorf("error creating RabbitMQClient: %s", err.Error())
        return
    }
    defer consumer.Close()

    messagesCh, err := consumer.Consume()
    if err != nil {
        t.Errorf("error executing consumer.Consumer: %s", err.Error())
        return
    }

    randomNum := rand.New(rand.NewSource(time.Now().Unix())).Int63()
    message := fmt.Sprintf("a message with a random number %d", randomNum)

    publisher, err := NewRabbitMQClient(
        WithQueueName("test-queue"),
        WithExchangeName("test-exchange"),
        WithExchangeType("topic"),
        WithConsumerRoutingKeys(testRoutingKey),
    )
    if err != nil {
        t.Error(err.Error())
        return
    }
    defer publisher.Close()

    n := 1000
    for i := 0; i < n; i++ {
        err = publisher.Publish(ctx, []byte(message), testRoutingKey)
        if err != nil {
            t.Errorf("error publishing message #%d: %s", i, err.Error())
            return
        }
        time.Sleep(1 * time.Millisecond)
    }

    waitForNMessagesWithTimeout(t, n, messagesCh, 5*time.Second)
}

func runRabbitMQ(ctx context.Context) func(ctx context.Context) {
    req := testcontainers.ContainerRequest{
        Image: "rabbitmq:3.8.0-management",
        Name:  "rabbitmq",
        User:  "rabbitmq",
        ExposedPorts: []string{
            fmt.Sprintf("%d:%d", DefaultPort, DefaultPort),
            fmt.Sprintf("%d:%d", ManagementUIPort, ManagementUIPort),
        },
        WaitingFor: wait.NewExecStrategy([]string{"rabbitmqadmin", "list", "queues"}).WithStartupTimeout(20 * time.Second),
    }
    mockserverC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        fmt.Printf("failed starting rabbitmq container: %s\n", err.Error())
        os.Exit(1)
    }

    return func(ctx context.Context) {
        if err := mockserverC.Terminate(ctx); err != nil {
            fmt.Printf("failed to terminate container: %s\n", err.Error())
        }
    }
}

func waitForMessageWithTimeout(t *testing.T, message string, messagesCh <-chan Message, duration time.Duration) {
    timeout := time.After(duration)
    select {
    case receivedMessage, ok := <-messagesCh:
        if !ok {
            t.Errorf("channel closed")
            return
        }
        assert.Equal(t, message, string(receivedMessage.Payload))
    case <-timeout:
        t.Errorf("timeout waiting message")
    }
}

func waitForNMessagesWithTimeout(t *testing.T, n int, messagesCh <-chan Message, duration time.Duration) {
    timeout := time.After(duration)
    messagesCount := 0
    for {
        select {
        case _, ok := <-messagesCh:
            if !ok {
                t.Errorf("channel closed")
                return
            }
            messagesCount++
            if messagesCount == n {
                return
            }
        case <-timeout:
            t.Errorf("timeout waiting messages")
            return
        }
    }
}
