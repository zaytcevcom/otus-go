package internalhttp

import (
	"context"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	server *http.Server
	logger Logger
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Application interface {
	GetEventsByDay(ctx context.Context, time time.Time) []storage.Event
	GetEventsByWeek(ctx context.Context, time time.Time) []storage.Event
	GetEventsByMonth(ctx context.Context, time time.Time) []storage.Event
	CreateEvent(ctx context.Context, title string, timeFrom time.Time, timeTo time.Time, description *string, userID string, notificationTime *time.Duration) (string, error)
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
}

func NewServer(logger Logger, app Application, host string, port int) *Server {

	server := &http.Server{
		Addr:         net.JoinHostPort(host, strconv.Itoa(port)),
		Handler:      loggingMiddleware(logger, NewHandler(logger, app)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		server: server,
		logger: logger,
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
