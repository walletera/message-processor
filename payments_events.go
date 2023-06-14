package main

import "encoding/json"

type EventEnvelope struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

type WithdrawalCreated struct {
    WithdrawalId       string  `json:"withdrawal_id"`
    Amount             float64 `json:"amount"`
    SourceAccount      string  `json:"source_account"`
    DestinationAccount string  `json:"destination_account"`
}
