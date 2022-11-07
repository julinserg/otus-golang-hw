package app_calendar

import (
	"database/sql"
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type AppCalendar struct {
	logger  app.Logger
	storage app.Storage
}

func eventFromStorage(event *storage.Event) app.Event {
	evApp := app.Event{
		ID: event.ID, Title: event.Title, Description: event.Description.String,
		UserID: event.UserID, NotificationTime: time.Duration(event.NotificationTime.Int64), TimeStart: event.TimeStart,
		TimeEnd: event.TimeEnd,
	}
	return evApp
}

func eventToStorage(event *app.Event) storage.Event {
	evStor := storage.Event{
		ID: event.ID, Title: event.Title,
		UserID: event.UserID, TimeStart: event.TimeStart,
		TimeEnd: event.TimeEnd,
	}

	descriptionIsValid := true
	if len(event.Description) == 0 {
		descriptionIsValid = false
	}
	NotificationTimeIsValid := true
	if event.NotificationTime == 0 {
		NotificationTimeIsValid = false
	}
	evStor.Description = sql.NullString{String: event.Description, Valid: descriptionIsValid}
	evStor.NotificationTime = sql.NullInt64{Int64: int64(event.NotificationTime), Valid: NotificationTimeIsValid}
	return evStor
}

func New(logger app.Logger, storage app.Storage) *AppCalendar {
	return &AppCalendar{logger, storage}
}

func (a *AppCalendar) AddEvent(event *app.Event) error {
	return a.storage.Add(eventToStorage(event))
}

func (a *AppCalendar) RemoveEvent(ID string) error {
	return a.storage.Remove(ID)
}

func (a *AppCalendar) UpdateEvent(event *app.Event) error {
	return a.storage.Update(eventToStorage(event))
}

type getEvent func(date time.Time) ([]storage.Event, error)

func (a *AppCalendar) genericGetEventsBy(date time.Time, f getEvent) ([]app.Event, error) {
	events, err := f(date)
	if err != nil {
		return nil, err
	}
	eventsApp := make([]app.Event, 0, len(events))
	for _, event := range events {
		eventsApp = append(eventsApp, eventFromStorage(&event))
	}
	return eventsApp, nil
}

func (a *AppCalendar) GetEventsByDay(date time.Time) ([]app.Event, error) {
	return a.genericGetEventsBy(date, func(date time.Time) ([]storage.Event, error) { return a.storage.GetEventsByDay(date) })
}

func (a *AppCalendar) GetEventsByMonth(date time.Time) ([]app.Event, error) {
	return a.genericGetEventsBy(date, func(date time.Time) ([]storage.Event, error) { return a.storage.GetEventsByMonth(date) })
}

func (a *AppCalendar) GetEventsByWeek(date time.Time) ([]app.Event, error) {
	return a.genericGetEventsBy(date, func(date time.Time) ([]storage.Event, error) { return a.storage.GetEventsByWeek(date) })
}
