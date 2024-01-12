package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID               string
	Title            string
	TimeFrom         time.Time `db:"time_from"`
	TimeTo           time.Time `db:"time_to"`
	Description      *string
	UserID           string         `db:"user_id"`
	NotificationTime *time.Duration `db:"notification_time"`
}

var (
	ErrDateBusy      = errors.New("this date is already booked by another event")
	ErrEventNotFound = errors.New("event not found")
)
