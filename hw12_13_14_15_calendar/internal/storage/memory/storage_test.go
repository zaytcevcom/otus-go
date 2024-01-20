package memorystorage

import (
	"context"
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

	ctx := context.Background()

	t.Run("CreateEvent", func(t *testing.T) {
		err := s.CreateEvent(ctx, event1)

		assert.NoError(t, err)
	})

	t.Run("CreateEventDuplicate", func(t *testing.T) {
		err := s.CreateEvent(ctx, event1)

		assert.ErrorIs(t, err, storage.ErrDateBusy)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		err := s.DeleteEvent(ctx, event1.ID)

		assert.Nil(t, err)
	})

	t.Run("DeleteUnknown", func(t *testing.T) {
		err := s.DeleteEvent(ctx, "unknownID")

		assert.ErrorIs(t, err, storage.ErrEventNotFound)
	})
}
