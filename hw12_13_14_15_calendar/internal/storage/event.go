package storage

import (
	"errors"
	"time"
)

var (
	ErrEventIdNotSet   = errors.New("Event ID not set")
	ErrEventIdNotExist = errors.New("Event ID not exist")
)

type Event struct {
	ID               string
	Title            string
	TimeStart        time.Time
	TimeEnd          time.Time
	Description      string
	UserID           string
	NotificationTime time.Duration
}
