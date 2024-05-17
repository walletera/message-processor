package webhook

import (
    "bytes"
    "context"
    "log/slog"
    "net/http"
    "os"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
    var tests = []struct {
        name                string
        httpMethod          string
        notificationPayload []byte
        expectedStatusCode  int
    }{
        {
            name:                "process request successfully",
            httpMethod:          "POST",
            notificationPayload: []byte("some payload"),
            expectedStatusCode:  http.StatusCreated,
        },
        {
            name:                "wrong http method produce MethodNotAllowed error",
            httpMethod:          "GET",
            notificationPayload: []byte("some payload"),
            expectedStatusCode:  http.StatusMethodNotAllowed,
        },
        {
            name:                "empty body produce StatusBadRequest",
            httpMethod:          "POST",
            notificationPayload: nil,
            expectedStatusCode:  http.StatusBadRequest,
        },
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            funcName(t, test.httpMethod, test.notificationPayload, test.expectedStatusCode)
        })
    }
}

func funcName(t *testing.T, httpMethod string, requestBody []byte, expectedStatusCode int) {
    server := NewServer(8282, ServerOpts{
        logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
    })
    msgCh, err := server.Consume()
    require.NoError(t, err)

    ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
    go func() {
        for msg := range msgCh {
            assert.NotNil(t, msg.Payload)
        }
        cancelCtx()
    }()

    httpReq, err := http.NewRequest(httpMethod, "http://127.0.0.1:8282/webhooks", bytes.NewReader(requestBody))
    require.NoError(t, err)

    resp, err := http.DefaultClient.Do(httpReq)
    require.NoError(t, err)

    assert.Equal(t, expectedStatusCode, resp.StatusCode)

    err = server.Close()
    require.NoError(t, err)

    <-ctx.Done()
    err = ctx.Err()
    assert.ErrorIs(t, err, context.Canceled)
}
