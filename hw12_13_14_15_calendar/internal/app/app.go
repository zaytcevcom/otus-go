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
	CreateEvent(event storage.Event) error
	UpdateEvent(id string, event storage.Event) error
	DeleteEvent(id string) error
	GetEventsByDay(time time.Time) []storage.Event
	GetEventsByWeek(time time.Time) []storage.Event
	GetEventsByMonth(time time.Time) []storage.Event
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

// Зачем тут нужен контекст?
func (a *App) CreateEvent(_ context.Context, id, title string) (string, error) {
	err := a.storage.CreateEvent(storage.Event{ID: id, Title: title})

	if errors.Is(err, storage.ErrDateBusy) {
		a.logger.Debug(storage.ErrDateBusy.Error())
		return "", storage.ErrDateBusy
	} else if err != nil {
		return "", err
	}

	a.logger.Debug("Created event: " + id)

	return id, nil
}
