package app

import (
	"context"
	"errors"
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
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEventsByDay(ctx context.Context, time time.Time) []storage.Event
	GetEventsByWeek(ctx context.Context, time time.Time) []storage.Event
	GetEventsByMonth(ctx context.Context, time time.Time) []storage.Event
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) (string, error) {
	err := a.storage.CreateEvent(ctx, storage.Event{ID: id, Title: title})

	if errors.Is(err, storage.ErrDateBusy) {
		return "", storage.ErrDateBusy
	} else if err != nil {
		return "", err
	}

	a.logger.Debug("Created event: " + id)

	return id, nil
}
