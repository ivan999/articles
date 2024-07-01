package response

import (
    "log/slog"
    "net/http"
    "encoding/json"
)

func renderJSON(
    w http.ResponseWriter, logger requestLogger, status int, value any,
) {
    jsonData, err := json.Marshal(value)
    if err != nil {
        const message = "failed rendering json"
        logger(slog.Error, message, "error", err, "value", value)
        http.Error(w, message, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(jsonData)
}

type responseData struct {
    Message string `json:"message"`
    Data    any    `json:"data"`
}

type responseDetails struct {
    Message string         `json:"message"`
    Details map[string]any `json:"details"`
}

type responseError struct {
    Error responseDetails `json:"error"`
}
