package main

import (
    "encoding/json"
    "fmt"
    "log"
)

type PaymentsEventsDeserializer struct {
}

func NewPaymentsEventsDeserializer() *PaymentsEventsDeserializer {
    return &PaymentsEventsDeserializer{}
}

func (d *PaymentsEventsDeserializer) Deserialize(rawPayload []byte) (Event[PaymentsEventsVisitor], error) {
    var event EventEnvelope
    err := json.Unmarshal(rawPayload, &event)
    if err != nil {
        return nil, fmt.Errorf("error deserializing message with payload %s: %w", rawPayload, err)
    }
    switch event.Type {
    case "WithdrawalCreated":
        var withdrawalCreated WithdrawalCreated
        err := json.Unmarshal(event.Data, &withdrawalCreated)
        if err != nil {
            log.Printf("error deserializing WithdrawalCreated event data %s: %s", event.Data, err.Error())
        }
        return withdrawalCreated, nil
    default:
        log.Printf("unexpected event type: %s", event.Type)
        return nil, nil
    }
}
