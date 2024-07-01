package api

import (
    "net/http"
    "context"
    "strings"

    "github.com/ivan999/articles/response"
)

func (usage *ServerUsage) authHandler(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        h := response.NewResponseHandler(w, r)
        details := map[string]any{}

        header := r.Header.Get("Authorization")
        if header == "" {
            const message = "empty authorization header"
            h.HandleClientError(http.StatusBadRequest, message, nil)
            return
        }

        headerParts := strings.Split(header, " ")
        if len(headerParts) != 2 {
            const message = "invalid authorization header format" 
            details[keyHeader] = header
            h.HandleClientError(http.StatusBadRequest, message, details)
            return
        }

        userID, err := parseJWT(headerParts[1])
        if err != nil {
            const message = "failed parsing jwt token"
            details[keyToken] = headerParts[1]
            details[keyError] = err.Error()
            h.HandleClientError(http.StatusBadRequest, message, details)
            return
        }

        ctx := context.WithValue(r.Context(), keyUserID, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}
