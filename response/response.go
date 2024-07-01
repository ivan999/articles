package response

import (
	"log/slog"
	"net/http"
)

type ResponseHandler struct {
    logger requestLogger
    writer http.ResponseWriter
}

func NewResponseHandler(
    w http.ResponseWriter, r *http.Request,
) *ResponseHandler {
    return &ResponseHandler{newRequestLogger(r), w}
}

func (h *ResponseHandler) HandleResponseData(
    status int, msg string, data any,
) {
    h.logger(slog.Info, msg)
    value := responseData{msg, data}
    renderJSON(h.writer, h.logger, status, value)
}

func (h *ResponseHandler) logDetails(
    method loggerMethod, msg string, details map[string]any,
) {
    args := make([]any, 0, len(details)*2)
    for key, value := range details {
        args = append(args, key, value)
    }
    h.logger(method, msg, args...)
}

func (h *ResponseHandler) HandleResponseDetails(
    status int, msg string, details map[string]any,
) {
    h.logDetails(slog.Info, msg, details)
    value := responseDetails{msg, details}
    renderJSON(h.writer, h.logger, status, value)
}

func (h *ResponseHandler) HandleClientError(
    status int, msg string, details map[string]any,
) {
    h.logDetails(slog.Warn, msg, details) 
    value := responseError{responseDetails{msg, details}}
    renderJSON(h.writer, h.logger, status, value)
}

func (h *ResponseHandler) HandleServerError(msg string, err error) {
    h.logger(slog.Error, msg, "error", err.Error())
    value := responseError{responseDetails{msg, nil}}
    renderJSON(h.writer, h.logger, http.StatusInternalServerError, value)
}
