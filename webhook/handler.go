package webhook

import (
    "context"
    "errors"
    "io"
    "log/slog"
    "net/http"
    "time"

    "github.com/walletera/message-processor/messages"
)

const (
    MessageProcessingTimeout = 10 * time.Second
)

type handler struct {
    logger *slog.Logger
    msgCh  chan messages.Message
}

func newHandler(msgCh chan messages.Message, logger *slog.Logger) *handler {
    return &handler{logger: logger, msgCh: msgCh}
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    if request.Method != http.MethodPost {
        writer.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    rawBody, err := io.ReadAll(request.Body)
    if err != nil {
        h.logError("failed reading request body", err)
        writer.WriteHeader(http.StatusInternalServerError)
        return
    }
    if rawBody == nil || len(rawBody) == 0 {
        h.logError("webhook request contains empty body", nil)
        writer.WriteHeader(http.StatusBadRequest)
        return
    }
    ctx, cancel := context.WithTimeout(context.Background(), MessageProcessingTimeout)
    defer cancel()
    acknowledger := NewAcknowledger(writer)
    h.msgCh <- messages.Message{
        Ctx:          ctx,
        Payload:      rawBody,
        Acknowledger: acknowledger,
    }
    select {
    case <-acknowledger.Done():
        return
    case <-ctx.Done():
        if errors.Is(ctx.Err(), context.DeadlineExceeded) {
            h.logError("timeout waiting for message to be processed", nil)
            writer.WriteHeader(http.StatusInternalServerError)
            return
        }
        return
    }
}

func (h *handler) logError(msg string, err error) {
    if h.logger == nil {
        return
    }
    if err != nil {
        h.logger.Error(msg, slog.String("error", err.Error()))
    } else {
        h.logger.Error(msg)
    }
}

func (h *handler) logDebug(msg string) {
    if h.logger == nil {
        return
    }
    h.logger.Debug(msg)
}
