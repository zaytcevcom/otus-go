package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(logger Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()

		next.ServeHTTP(w, r)

		latency := time.Since(startTime)

		logger.Info(
			fmt.Sprintf(
				"%s [%s] %s %s HTTP/%s %d %s \"%s\"",
				r.RemoteAddr,
				time.Now().Format(time.RFC1123Z),
				r.Method, r.URL.Path, r.Proto,
				http.StatusOK, // можно узнать из response ???
				latency,
				r.Header.Get("User-Agent"),
			),
		)
	})
}
