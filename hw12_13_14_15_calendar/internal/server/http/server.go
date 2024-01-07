package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	server *http.Server
	app    Application
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

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello!"))
	if err != nil {
		return
	}
}

func NewServer(logger Logger, app Application, host string, port int) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", HelloHandler{
		logger: logger,
	})

	server := &http.Server{
		Addr:         net.JoinHostPort(host, strconv.Itoa(port)),
		Handler:      loggingMiddleware(logger, mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		server: server,
		app:    app,
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
