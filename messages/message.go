package messages

import "context"

type Message struct {
    Ctx          context.Context
    Payload      []byte
    Acknowledger Acknowledger
}
