package storage

import "time"

type Notification struct {
	ID       string
	Title    string
	Datetime time.Time
	UserID   string
}
