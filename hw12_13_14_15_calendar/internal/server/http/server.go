package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	server *http.Server
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Application interface{}

type HelloHandler struct {
	logger Logger
}

func (h HelloHandler) ServeHTTP(_ http.ResponseWriter, r *http.Request) {
	latency, ok := r.Context().Value(LatencyKey).(time.Duration)

	if !ok {
		latency = 0
	}

	h.logger.Info(
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
}

// Зачем нужен Application?
// Зачем Logger, если в main.go в calendar уже содержит логгер? // calendar := app.New(logg, storage).
func NewServer(logger Logger, _ Application, host string, port int) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", HelloHandler{
		logger: logger,
	})

	server := &http.Server{
		Addr:         net.JoinHostPort(host, strconv.Itoa(port)),
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		server: server,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
