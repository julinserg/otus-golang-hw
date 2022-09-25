package app

import (
	"context"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface { // TODO
	Get(id string) storage.Event
	Add(event storage.Event) error
	Update(event storage.Event) error
	Remove(id string) error
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
