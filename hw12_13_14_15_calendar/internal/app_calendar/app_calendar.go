package app_calendar

import (
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type AppCalendar struct {
	logger  app.Logger
	storage app.Storage
}

func New(logger app.Logger, storage app.Storage) *AppCalendar {
	return &AppCalendar{logger, storage}
}

func (a *AppCalendar) AddEvent(event *app.EventApp) error {
	return a.storage.Add(storage.Event{ID: event.ID, Title: event.Title, Description: event.Description,
		UserID: event.UserID, NotificationTime: event.NotificationTime, TimeStart: event.TimeStart,
		TimeEnd: event.TimeEnd})
}

func (a *AppCalendar) RemoveEvent(ID string) error {
	return a.storage.Remove(ID)
}

func (a *AppCalendar) UpdateEvent(event *app.EventApp) error {
	return a.storage.Update(storage.Event{ID: event.ID, Title: event.Title, Description: event.Description,
		UserID: event.UserID, NotificationTime: event.NotificationTime, TimeStart: event.TimeStart,
		TimeEnd: event.TimeEnd})
}

type getEvent func(date time.Time) ([]storage.Event, error)

func (a *AppCalendar) genericGetEventsBy(date time.Time, f getEvent) ([]app.EventApp, error) {
	events, err := f(date)
	if err != nil {
		return nil, err
	}
	eventsApp := make([]app.EventApp, 0, len(events))
	for _, event := range events {
		eventsApp = append(eventsApp, app.EventApp{ID: event.ID, Title: event.Title, Description: event.Description,
			UserID: event.UserID, NotificationTime: event.NotificationTime, TimeStart: event.TimeStart,
			TimeEnd: event.TimeEnd})
	}
	return eventsApp, nil
}

func (a *AppCalendar) GetEventsByDay(date time.Time) ([]app.EventApp, error) {
	return a.genericGetEventsBy(date, func(date time.Time) ([]storage.Event, error) { return a.storage.GetEventsByDay(date) })
}

func (a *AppCalendar) GetEventsByMonth(date time.Time) ([]app.EventApp, error) {
	return a.genericGetEventsBy(date, func(date time.Time) ([]storage.Event, error) { return a.storage.GetEventsByMonth(date) })
}

func (a *AppCalendar) GetEventsByWeek(date time.Time) ([]app.EventApp, error) {
	return a.genericGetEventsBy(date, func(date time.Time) ([]storage.Event, error) { return a.storage.GetEventsByWeek(date) })
}
