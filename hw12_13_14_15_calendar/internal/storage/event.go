package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID               string
	Title            string
	TimeFrom         time.Time
	TimeTo           time.Time
	Description      *string
	UserID           string
	NotificationTime *time.Duration
}

var (
	ErrDateBusy      = errors.New("the date is already booked by another event")
	ErrEventNotFound = errors.New("event not found")
)
