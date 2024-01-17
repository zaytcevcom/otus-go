package app

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zaytcevcom/otus-go/hw12_13_14_15_calendar/internal/storage"
	"testing"
	"time"
)

type MockStorage struct {
	mock.Mock
}

func (ms *MockStorage) GetEventsByDay(ctx context.Context, time time.Time) []storage.Event {
	args := ms.Called(ctx, time)
	return args.Get(0).([]storage.Event)
}

func (ms *MockStorage) GetEventsByWeek(ctx context.Context, time time.Time) []storage.Event {
	args := ms.Called(ctx, time)
	return args.Get(0).([]storage.Event)
}

func (ms *MockStorage) GetEventsByMonth(ctx context.Context, time time.Time) []storage.Event {
	args := ms.Called(ctx, time)
	return args.Get(0).([]storage.Event)
}

func (ms *MockStorage) CreateEvent(ctx context.Context, event storage.Event) (string, error) {
	args := ms.Called(ctx, event)
	return args.String(0), args.Error(1)
}

func (ms *MockStorage) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	args := ms.Called(ctx, id, event)
	return args.Error(0)
}

func (ms *MockStorage) DeleteEvent(ctx context.Context, id string) error {
	args := ms.Called(ctx, id)
	return args.Error(0)
}

func (ms *MockStorage) Connect(ctx context.Context) error {
	args := ms.Called(ctx)
	return args.Error(0)
}

func (ms *MockStorage) Close(ctx context.Context) error {
	args := ms.Called(ctx)
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (ml *MockLogger) Debug(msg string) {
	ml.Called(msg)
}

func (ml *MockLogger) Info(msg string) {
	ml.Called(msg)
}

func (ml *MockLogger) Warn(msg string) {
	ml.Called(msg)
}

func (ml *MockLogger) Error(msg string) {
	ml.Called(msg)
}

func TestApp(t *testing.T) {
	mockStorage := new(MockStorage)
	mockLogger := new(MockLogger)
	app := New(mockLogger, mockStorage)

	event := storage.Event{
		ID:               uuid.NewString(),
		Title:            "Test Event",
		TimeFrom:         time.Now(),
		TimeTo:           time.Now().Add(1 * time.Hour),
		Description:      nil,
		UserID:           "1",
		NotificationTime: nil,
	}

	t.Run("GetEventsByDay", func(t *testing.T) {
		mockStorage.On("GetEventsByDay", mock.Anything, mock.Anything).Return([]storage.Event{event})
		app.GetEventsByDay(context.Background(), time.Now())
		mockStorage.AssertExpectations(t)
	})

	t.Run("GetEventsByWeek", func(t *testing.T) {
		mockStorage.On("GetEventsByWeek", mock.Anything, mock.Anything).Return([]storage.Event{event})
		app.GetEventsByWeek(context.Background(), time.Now())
		mockStorage.AssertExpectations(t)
	})

	t.Run("GetEventsByMonth", func(t *testing.T) {
		mockStorage.On("GetEventsByMonth", mock.Anything, mock.Anything).Return([]storage.Event{event})
		app.GetEventsByMonth(context.Background(), time.Now())
		mockStorage.AssertExpectations(t)
	})

	t.Run("CreateEvent", func(t *testing.T) {
		mockStorage.On("CreateEvent", mock.Anything, mock.Anything).Return(event.ID, nil)
		mockLogger.On("Debug", mock.AnythingOfType("string")).Return()
		_, err := app.CreateEvent(context.Background(), event.Title, event.TimeFrom, event.TimeTo, event.Description, event.UserID, event.NotificationTime)
		if err != nil {
			assert.Fail(t, err.Error())
		}
		mockStorage.AssertExpectations(t)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		mockStorage.On("UpdateEvent", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		err := app.UpdateEvent(context.Background(), "111111111111", event)
		if err != nil {
			assert.Fail(t, err.Error())
		}
		mockStorage.AssertExpectations(t)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		mockStorage.On("DeleteEvent", mock.Anything, mock.Anything).Return(nil)
		err := app.DeleteEvent(context.Background(), "111111111111")
		if err != nil {
			assert.Fail(t, err.Error())
		}
		mockStorage.AssertExpectations(t)
	})
}
