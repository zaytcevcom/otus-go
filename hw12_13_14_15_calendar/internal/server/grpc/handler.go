package internalgrpc

import (
	"context"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/api/proto"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Handler struct {
	proto.UnimplementedEventServiceServer
	logger Logger
	app    Application
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewHandler(logger Logger, app Application) *Handler {
	return &Handler{
		logger: logger,
		app:    app,
	}
}

func (s *Handler) GetEventsByDay(ctx context.Context, r *proto.GetEventsRequest) (*proto.GetEventsResponse, error) {

	events := s.app.GetEventsByDay(ctx, time.Unix(r.Time, 0))

	eventsResponse := &proto.GetEventsResponse{
		Events: make([]*proto.EventResponse, len(events)),
	}

	for _, event := range events {
		eventResponse := &proto.EventResponse{
			Id:               event.ID,
			Title:            event.Title,
			TimeFrom:         event.TimeFrom.Unix(),
			TimeTo:           event.TimeTo.Unix(),
			Description:      *event.Description,
			UserId:           event.UserID,
			NotificationTime: event.NotificationTime.Milliseconds(),
		}
		eventsResponse.Events = append(eventsResponse.Events, eventResponse)
	}

	return eventsResponse, nil
}

func (s *Handler) GetEventsByWeek(ctx context.Context, r *proto.GetEventsRequest) (*proto.GetEventsResponse, error) {

	events := s.app.GetEventsByWeek(ctx, time.Unix(r.Time, 0))

	eventsResponse := &proto.GetEventsResponse{
		Events: make([]*proto.EventResponse, len(events)),
	}

	for _, event := range events {
		eventResponse := &proto.EventResponse{
			Id:               event.ID,
			Title:            event.Title,
			TimeFrom:         event.TimeFrom.Unix(),
			TimeTo:           event.TimeTo.Unix(),
			Description:      *event.Description,
			UserId:           event.UserID,
			NotificationTime: event.NotificationTime.Milliseconds(),
		}
		eventsResponse.Events = append(eventsResponse.Events, eventResponse)
	}

	return eventsResponse, nil
}

func (s *Handler) GetEventsByMonth(ctx context.Context, r *proto.GetEventsRequest) (*proto.GetEventsResponse, error) {

	events := s.app.GetEventsByMonth(ctx, time.Unix(r.Time, 0))

	eventsResponse := &proto.GetEventsResponse{
		Events: make([]*proto.EventResponse, len(events)),
	}

	for _, event := range events {
		eventResponse := &proto.EventResponse{
			Id:               event.ID,
			Title:            event.Title,
			TimeFrom:         event.TimeFrom.Unix(),
			TimeTo:           event.TimeTo.Unix(),
			Description:      *event.Description,
			UserId:           event.UserID,
			NotificationTime: int64(event.NotificationTime.Seconds()),
		}
		eventsResponse.Events = append(eventsResponse.Events, eventResponse)
	}

	return eventsResponse, nil
}

func (s *Handler) CreateEvent(ctx context.Context, r *proto.CreateEventRequest) (*proto.EventResponse, error) {

	duration := time.Duration(r.NotificationTime) * time.Second

	id, err := s.app.CreateEvent(
		ctx,
		r.Title,
		time.Unix(r.TimeFrom, 0),
		time.Unix(r.TimeTo, 0),
		&r.Description,
		r.UserId,
		&duration,
	)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &proto.EventResponse{
		Id:               id,
		Title:            r.Title,
		TimeFrom:         r.TimeFrom,
		TimeTo:           r.TimeTo,
		Description:      r.Description,
		UserId:           r.UserId,
		NotificationTime: r.NotificationTime,
	}, nil
}

func (s *Handler) UpdateEvent(ctx context.Context, r *proto.UpdateEventRequest) (*proto.EventResponse, error) {

	duration := time.Duration(r.NotificationTime) * time.Second

	event := storage.Event{
		ID:               r.Id,
		Title:            r.Title,
		TimeFrom:         time.Unix(r.TimeFrom, 0),
		TimeTo:           time.Unix(r.TimeTo, 0),
		Description:      &r.Description,
		UserID:           r.UserId,
		NotificationTime: &duration,
	}

	err := s.app.UpdateEvent(ctx, r.Id, event)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &proto.EventResponse{
		Id:               r.Id,
		Title:            r.Title,
		TimeFrom:         r.TimeFrom,
		TimeTo:           r.TimeTo,
		Description:      r.Description,
		UserId:           r.UserId,
		NotificationTime: r.NotificationTime,
	}, nil
}

func (s *Handler) DeleteEvent(ctx context.Context, r *proto.DeleteEventRequest) (*proto.DeleteEventResponse, error) {

	err := s.app.DeleteEvent(ctx, r.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &proto.DeleteEventResponse{}, nil
}
