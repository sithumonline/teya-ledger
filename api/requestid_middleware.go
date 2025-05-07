package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const xRequestIDHeaderKey = "X-Request-Id"

// RequestIDMiddleware for tracing purpose
// This will be handy with telemetry and observability
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(xRequestIDHeaderKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set(xRequestIDHeaderKey, requestID)
		ctx := context.WithValue(r.Context(), xRequestIDHeaderKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
