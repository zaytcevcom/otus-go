package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) GetEventsByDay(ctx context.Context, t time.Time) []storage.Event {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 0, 1)

	return s.getEventsByPeriod(start, end)
}

func (s *Storage) GetEventsByWeek(ctx context.Context, t time.Time) []storage.Event {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 0, 7)

	return s.getEventsByPeriod(start, end)
}

func (s *Storage) GetEventsByMonth(ctx context.Context, t time.Time) []storage.Event {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := start.AddDate(0, 1, 0)

	return s.getEventsByPeriod(start, end)
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isEventTimeBusy(event.TimeFrom, "") {
		return "", storage.ErrDateBusy
	}

	s.events[event.ID] = event

	return event.ID, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.events[id]
	if !exists {
		return storage.ErrEventNotFound
	}

	if s.isEventTimeBusy(event.TimeFrom, id) {
		return storage.ErrDateBusy
	}

	s.events[id] = event

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.events[id]
	if !exists {
		return storage.ErrEventNotFound
	}
	delete(s.events, id)

	return nil
}

func (s *Storage) getEventsByPeriod(start time.Time, end time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []storage.Event

	for _, event := range s.events {
		//if event.TimeFrom.After(start) && event.TimeFrom.Before(end) {
		events = append(events, event)
		//}
	}

	return events
}

func (s *Storage) isEventTimeBusy(t time.Time, eventID string) bool {
	for id, e := range s.events {
		if id != eventID && e.TimeFrom.Equal(t) {
			return true
		}
	}
	return false
}
