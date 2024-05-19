package messages

type ErrorCode int

const (
    UnprocessableMessage ErrorCode = iota + 1
    InternalError
)

type NackOpts struct {
    Requeue      bool
    Code         ErrorCode
    ErrorMessage string
}

type Acknowledger interface {
    Ack() error
    Nack(opts NackOpts) error
}
