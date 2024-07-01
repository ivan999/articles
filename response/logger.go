package response

import (
    "net/http"
)

type loggerMethod func(string, ...any)

type requestLogger func(loggerMethod, string, ...any)

func newRequestLogger(r *http.Request) requestLogger {
    return func(method loggerMethod, msg string, args ...any) {
        requestArgs := []any{"Method", r.Method, "URI", r.RequestURI}
        method(msg, append(requestArgs, args...)...)
    }
}
