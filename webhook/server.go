package webhook

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "net"
    "net/http"
    "time"

    "github.com/walletera/message-processor/messages"
)

const (
    DefaultPort = 8282

    shutdownTimeout = 10 * time.Second
)

type Server struct {
    httpServer http.Server
    msgCh      chan messages.Message
    logger     *slog.Logger
}

type ServerOpts struct {
    httpHandler http.Handler
    logger      *slog.Logger
}

func NewServer(port int, opts ServerOpts) *Server {
    msgCh := make(chan messages.Message)
    server := &Server{
        httpServer: http.Server{
            Addr:    fmt.Sprintf(":%d", port),
            Handler: newHandler(msgCh, opts.logger),
        },
        msgCh:  msgCh,
        logger: opts.logger,
    }
    return server
}

func (s *Server) Consume() (<-chan messages.Message, error) {
    listener, err := net.Listen("tcp", s.httpServer.Addr)
    if err != nil {
        return nil, fmt.Errorf("failed listening on %s: %w", s.httpServer.Addr, err)
    }
    go func() {
        if err := s.httpServer.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
            s.logError("http server error", err)
            close(s.msgCh)
        }
    }()

    return s.msgCh, nil
}

func (s *Server) Close() error {
    defer close(s.msgCh)

    shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), shutdownTimeout)
    defer shutdownRelease()

    if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
        return fmt.Errorf("http shutdown error: %w", err)
    }

    return nil
}

func (s *Server) logError(msg string, err error) {
    if s.logger == nil {
        return
    }
    s.logger.Error(msg, slog.String("error", err.Error()))
}
