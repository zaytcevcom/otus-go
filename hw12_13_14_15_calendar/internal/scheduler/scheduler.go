package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type Scheduler struct {
	logger  Logger
	storage Storage
	broker  MessageBroker
	doneCh  chan interface{}
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	GetNeedNotify(ctx context.Context, time time.Time) []storage.Event
	DeleteOldest(ctx context.Context, time time.Time) error
}

type MessageBroker interface {
	Publish(body string) error
	Close() error
}

func New(logger Logger, storage Storage, broker MessageBroker) *Scheduler {
	return &Scheduler{
		logger:  logger,
		storage: storage,
		broker:  broker,
		doneCh:  make(chan interface{}),
	}
}

func (s Scheduler) Start(ctx context.Context, interval int) error {

	s.logger.Debug("Scheduled started!")

	t := time.NewTicker(time.Duration(interval) * time.Second)

	for {
		select {
		case <-s.doneCh:
			return nil
		case <-t.C:

			s.logger.Debug("Tick!")

			timeNow := time.Now()
			yearAgo := timeNow.AddDate(-1, 0, 0)

			if err := s.deleteOldest(ctx, yearAgo); err != nil {
				s.logger.Error(err.Error())
			}

			if err := s.notifyEvents(ctx, timeNow); err != nil {
				s.logger.Error(err.Error())
			}
		}
	}
}

func (s Scheduler) Stop() error {
	close(s.doneCh)
	return s.broker.Close()
}

func (s Scheduler) notifyEvents(ctx context.Context, t time.Time) error {
	events := s.storage.GetNeedNotify(ctx, t)

	for _, event := range events {

		notification := storage.Notification{
			ID:       uuid.NewString(),
			Title:    event.Title,
			Datetime: event.TimeFrom,
			UserID:   event.UserID,
		}

		msg, err := json.Marshal(notification)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed marshal: %s", err))
			return err
		}

		err = s.broker.Publish(string(msg))
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed publish: %s", err))
			return err
		}
	}

	return nil
}

func (s Scheduler) deleteOldest(ctx context.Context, t time.Time) error {
	return s.storage.DeleteOldest(ctx, t)
}
