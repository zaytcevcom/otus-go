package internalhttp

import (
	"context"
	"net/http"
	"time"
)

type contextKey string

const LatencyKey = contextKey("Latency")

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		latency := time.Since(start)

		ctx := context.WithValue(r.Context(), LatencyKey, latency)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
