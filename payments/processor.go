package payments

import "github.com/walletera/message-processor/messages"

// NewProcessor returns a messages.Processor which is specific to payments events
func NewProcessor(messageConsumer messages.Consumer, eventsVisitor EventsVisitor, errorHandler messages.ErrorHandler) *messages.Processor[EventsVisitor] {
    return messages.NewProcessor[EventsVisitor](messageConsumer, NewEventsDeserializer(), eventsVisitor, errorHandler)
}
