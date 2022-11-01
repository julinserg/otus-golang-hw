package app

import (
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type EventApp struct {
	ID               string        `json:"id"`
	Title            string        `json:"title"`
	TimeStart        time.Time     `json:"time_start"`
	TimeEnd          time.Time     `json:"time_stop"`
	Description      string        `json:"description"`
	UserID           string        `json:"user_id"`
	NotificationTime time.Duration `json:"time_notify"`
}

type NotifyEvent struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	TimeStart time.Time `json:"time_start"`
	UserID    string    `json:"user_id"`
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type Storage interface {
	Add(event storage.Event) error
	Update(event storage.Event) error
	Remove(id string) error
	GetEventsByDay(date time.Time) ([]storage.Event, error)
	GetEventsByWeek(dateBeginWeek time.Time) ([]storage.Event, error)
	GetEventsByMonth(dateBeginMonth time.Time) ([]storage.Event, error)
}
