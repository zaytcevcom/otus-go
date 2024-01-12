package internalgrpc

import (
	"context"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/api/proto"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
	"time"
)

type Server struct {
	server *grpc.Server
	logger Logger
	app    Application
	host   string
	port   int
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

	return &Server{
		server: grpc.NewServer(UnaryInterceptor(logger)),
		logger: logger,
		app:    app,
		host:   host,
		port:   port,
	}
}

func (s *Server) Start(ctx context.Context) error {

	listener, err := net.Listen("tcp", net.JoinHostPort(s.host, strconv.Itoa(s.port)))
	if err != nil {
		return err
	}

	proto.RegisterEventServiceServer(s.server, NewHandler(s.logger, s.app))
	reflection.Register(s.server)

	err = s.server.Serve(listener)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
