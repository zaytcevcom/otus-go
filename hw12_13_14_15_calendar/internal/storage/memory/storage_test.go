package memorystorage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	s := New()

	timeFrom := time.Now()
	timeTo := timeFrom.Add(1 * time.Hour)
	event1 := storage.Event{
		ID:       "111111111",
		Title:    "Event1",
		TimeFrom: timeFrom,
		TimeTo:   timeTo,
	}

	t.Run("CreateEvent", func(t *testing.T) {
		err := s.CreateEvent(event1)

		assert.Nil(t, err)
	})

	t.Run("CreateEventDuplicate", func(t *testing.T) {
		err := s.CreateEvent(event1)

		assert.ErrorIs(t, err, storage.ErrDateBusy)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		err := s.DeleteEvent(event1.ID)

		assert.Nil(t, err)
	})

	t.Run("DeleteUnknown", func(t *testing.T) {
		err := s.DeleteEvent("unknownID")

		assert.ErrorIs(t, err, storage.ErrEventNotFound)
	})
}
