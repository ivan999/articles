package api

import (
    "io"
    "fmt"
    "errors"
    "strconv"
    "net/http"
    "encoding/json"

    "github.com/ivan999/articles/response"
    "github.com/ivan999/articles/storage"
)

type receiveError struct {
    internal bool
    message  string
    wrapped  error
}

func (rerr *receiveError) Error() string {
    return fmt.Sprintf("internal=%v message=%s wrapped=%v",
        rerr.internal, rerr.message, rerr.wrapped)
}

func receiveJSON(r *http.Request, data any) error {
    defer r.Body.Close()

    body, err := io.ReadAll(r.Body)
    if err != nil {
        const message = "failed reading request body"
        return &receiveError{true, message, err}
    }

    err = json.Unmarshal(body, &data)
    if err != nil {
        const message = "failed unmarshaling request body"
        return &receiveError{false, message, err}
    }

    return nil
}

func receiveParam(r *http.Request, key string) (int64, error) {
    value := r.URL.Query().Get(key)
    number, err := strconv.ParseInt(value, 10, 64)
    if err != nil {
        message := "failed parsing " + key + " parameter"
        return 0, &receiveError{false, message, err}
    }

    return number, nil
}

func handleReceiveError(h *response.ResponseHandler, err error) {
    var rerr *receiveError
    ok := errors.As(err, &rerr)
    if !ok {
        h.HandleServerError(
            "failed receive error handling",
            errors.New("invalid error type"),
        )
        return
    }

    if rerr.internal {
        h.HandleServerError(rerr.message, rerr.wrapped)
    } else {
        details := map[string]any{keyError: rerr.wrapped.Error()}
        h.HandleClientError(http.StatusBadRequest, rerr.message, details)
    }
}

func handleStorageError(
    h *response.ResponseHandler, err error, details map[string]any,
) {
    var serr *storage.StorageError
    ok := errors.As(err, &serr)
    if !ok {
        h.HandleServerError(
            "failed storage error handling", 
            errors.New("invalid error type"),
        )
        return
    }

    switch serr.Code {
    case storage.ErrCodeNotFound:
        h.HandleClientError(http.StatusNotFound, serr.Message, details)
    case storage.ErrCodeUniqueExists:
        h.HandleClientError(http.StatusConflict, serr.Message, details)
    default:
        h.HandleServerError(serr.Message, serr.Wrapped)
    }
}
