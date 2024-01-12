package app

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	GetEventsByDay(ctx context.Context, time time.Time) []storage.Event
	GetEventsByWeek(ctx context.Context, time time.Time) []storage.Event
	GetEventsByMonth(ctx context.Context, time time.Time) []storage.Event
	CreateEvent(ctx context.Context, event storage.Event) (string, error)
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) GetEventsByDay(ctx context.Context, time time.Time) []storage.Event {
	return a.storage.GetEventsByDay(ctx, time)
}

func (a *App) GetEventsByWeek(ctx context.Context, time time.Time) []storage.Event {
	return a.storage.GetEventsByWeek(ctx, time)
}

func (a *App) GetEventsByMonth(ctx context.Context, time time.Time) []storage.Event {
	return a.storage.GetEventsByMonth(ctx, time)
}

func (a *App) CreateEvent(
	ctx context.Context,
	title string,
	timeFrom time.Time,
	timeTo time.Time,
	description *string,
	userID string,
	notificationTime *time.Duration,
) (string, error) {

	event := storage.Event{
		ID:               uuid.NewString(),
		Title:            title,
		TimeFrom:         timeFrom,
		TimeTo:           timeTo,
		Description:      description,
		UserID:           userID,
		NotificationTime: notificationTime,
	}

	id, err := a.storage.CreateEvent(ctx, event)

	if err != nil {
		return "", err
	}

	a.logger.Debug("Created event: " + id)

	return id, nil
}

func (a *App) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, id, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}
