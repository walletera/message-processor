package main

import (
    "log"
)

func main() {

    processor := NewMessageProcessor()

    err := processor.Start()
    if err != nil {
        log.Fatalf("failed to start message processor: %s", err.Error())
    }

    blockForeverCh := make(chan any)
    <-blockForeverCh
}
