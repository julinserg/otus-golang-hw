package sqlstorage

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	var err error
	s.db, err = sqlx.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}
	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) Get(id string) storage.Event {
	return storage.Event{}
}

func (s *Storage) Add(event storage.Event) error {

	_, err := s.db.NamedExec(`INSERT INTO events (id,title,time_start,time_stop,description,user_id,time_notify)
	 VALUES (:id,:title,:time_start,:time_stop,:description,:user_id,:time_notify)`,
		map[string]interface{}{
			"id":          event.ID,
			"title":       event.Title,
			"time_start":  event.TimeStart,
			"time_stop":   event.TimeEnd,
			"description": event.Description,
			"user_id":     event.UserID,
			"time_notify": event.NotificationTime,
		})

	return err
}

func (s *Storage) Update(event storage.Event) error {
	_, err := s.db.NamedExec(`UPDATE events SET title=:title, time_start=:time_start,
	 time_stop=:time_stop,description=:description, user_id =:user_id, time_notify=:time_notify WHERE id = `+`'`+event.ID+`'`,
		map[string]interface{}{
			"title":       event.Title,
			"time_start":  event.TimeStart,
			"time_stop":   event.TimeEnd,
			"description": event.Description,
			"user_id":     event.UserID,
			"time_notify": event.NotificationTime,
		})
	return err
}

func (s *Storage) Remove(id string) error {
	return nil
}
