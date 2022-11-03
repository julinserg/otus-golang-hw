package sqlstorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

var schema = `
DROP table if exists events;
CREATE table events (
    id              text primary key,
    title           text not null,
    time_start      timestamp not null,
    time_stop       timestamp not null,
    description     text,
    user_id         text not null,    
    time_notify     bigint,
	is_notifyed     boolean,
	CONSTRAINT time_start_unique UNIQUE (time_start)
);`

func TestStorageBasic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := "host=localhost port=5432 user=sergey password=sergey dbname=calendar_test sslmode=disable"
	dbTest, err := sqlx.Open("pgx", dsn)
	require.Nil(t, err)

	err = dbTest.PingContext(ctx)
	require.Nil(t, err)

	dbTest.MustExec(schema)

	dbTest.Close()

	st := New()

	err = st.Connect(ctx, dsn)
	require.Nil(t, err)
	defer func() {
		if err := st.Close(ctx); err != nil {
			fmt.Printf("cannot close psql connection: " + err.Error())
		}
	}()

	event1 := storage.Event{
		ID: "{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}", Title: "event 1",
		TimeStart: time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC),
	}
	err = st.Add(event1)
	require.Nil(t, err)
	event2 := storage.Event{
		ID: "{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}", Title: "event 2",
		TimeStart: time.Date(2022, time.Month(2), 1, 1, 10, 30, 0, time.UTC),
	}
	err = st.Add(event2)
	require.Nil(t, err)

	resGet, _ := st.Get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}")
	require.Equal(t, event1, resGet)
	resGet, _ = st.Get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}")
	require.Equal(t, event2, resGet)
	resGet, _ = st.Get("{81e125ce-072e-4556-8a4c-597572a7277a}")
	require.Equal(t, storage.Event{}, resGet)

	resGet, _ = st.Get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}")
	require.Equal(t, "event 1", resGet.Title)
	resGet, _ = st.Get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}")
	require.Equal(t, "event 2", resGet.Title)

	event1.Title = "event 5"
	st.Update(event1)

	resGet, _ = st.Get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}")
	require.Equal(t, "event 5", resGet.Title)
	resGet, _ = st.Get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}")
	require.Equal(t, "event 2", resGet.Title)

	st.Remove(event1.ID)
	resGet, _ = st.Get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}")
	require.Equal(t, storage.Event{}, resGet)
	resGet, _ = st.Get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}")
	require.Equal(t, event2, resGet)

	event1.ID = "{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}"
	event1.Title = "event 6"
	require.ErrorIs(t, st.Update(event1), storage.ErrEventIDNotExist)

	require.ErrorIs(t, st.Remove(event1.ID), storage.ErrEventIDNotExist)

	event3 := storage.Event{Title: "event 3"}
	require.ErrorIs(t, st.Add(event3), storage.ErrEventIDNotSet)
}

func TestStorageLogic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := "host=localhost port=5432 user=sergey password=sergey dbname=calendar_test sslmode=disable"
	dbTest, err := sqlx.Open("pgx", dsn)
	require.Nil(t, err)

	err = dbTest.PingContext(ctx)
	require.Nil(t, err)

	dbTest.MustExec(schema)

	dbTest.Close()

	st := New()

	err = st.Connect(ctx, dsn)
	require.Nil(t, err)
	defer func() {
		if err := st.Close(ctx); err != nil {
			fmt.Printf("cannot close psql connection: " + err.Error())
		}
	}()

	st.Add(storage.Event{
		ID:        "1",
		Title:     "event 1",
		TimeStart: time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC),
	})

	st.Add(storage.Event{
		ID:        "2",
		Title:     "event 2",
		TimeStart: time.Date(2022, time.Month(2), 1, 1, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(2), 1, 1, 10, 30, 0, time.UTC),
	})

	st.Add(storage.Event{
		ID:        "3",
		Title:     "event 3",
		TimeStart: time.Date(2022, time.Month(1), 2, 1, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(1), 2, 1, 10, 30, 0, time.UTC),
	})

	st.Add(storage.Event{
		ID:        "4",
		Title:     "event 4",
		TimeStart: time.Date(2022, time.Month(1), 1, 2, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(1), 1, 2, 10, 30, 0, time.UTC),
	})

	st.Add(storage.Event{
		ID:        "5",
		Title:     "event 5",
		TimeStart: time.Date(2022, time.Month(1), 3, 2, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(1), 3, 2, 10, 30, 0, time.UTC),
	})

	st.Add(storage.Event{
		ID:        "6",
		Title:     "event 6",
		TimeStart: time.Date(2022, time.Month(1), 20, 2, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(1), 20, 2, 10, 30, 0, time.UTC),
	})

	st.Add(storage.Event{
		ID:        "7",
		Title:     "event 7",
		TimeStart: time.Date(2022, time.Month(3), 20, 2, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(3), 20, 2, 10, 30, 0, time.UTC),
	})

	st.Add(storage.Event{
		ID:        "8",
		Title:     "event 8",
		TimeStart: time.Date(2022, time.Month(4), 20, 2, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(4), 20, 2, 10, 30, 0, time.UTC),
	})

	res, err := st.GetEventsByDay(time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC))
	require.Nil(t, err)
	require.Equal(t, 2, len(res))

	res, err = st.GetEventsByWeek(time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC))
	require.Nil(t, err)
	require.Equal(t, 4, len(res))

	res, err = st.GetEventsByMonth(time.Date(2022, time.Month(2), 1, 1, 10, 30, 0, time.UTC))
	require.Nil(t, err)
	require.Equal(t, 1, len(res))

	res, err = st.GetEventsByMonth(time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC))
	require.Nil(t, err)
	require.Equal(t, 6, len(res))

	err = st.Add(storage.Event{
		ID:        "9",
		Title:     "event 9",
		TimeStart: time.Date(2022, time.Month(4), 20, 2, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(4), 20, 2, 10, 30, 0, time.UTC),
	})
	require.ErrorIs(t, err, storage.ErrTimeBusy)

	err = st.Update(storage.Event{
		ID:        "8",
		Title:     "event 9",
		TimeStart: time.Date(2022, time.Month(3), 20, 2, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(3), 20, 2, 10, 30, 0, time.UTC),
	})
	require.ErrorIs(t, err, storage.ErrTimeBusy)

	st.Remove("8")

	err = st.Add(storage.Event{
		ID:        "9",
		Title:     "event 9",
		TimeStart: time.Date(2022, time.Month(4), 20, 2, 10, 30, 0, time.UTC),
		TimeEnd:   time.Date(2022, time.Month(4), 20, 2, 10, 30, 0, time.UTC),
	})
	require.Nil(t, err)
}

func TestStorageGetEventForNotify(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := "host=localhost port=5432 user=sergey password=sergey dbname=calendar_test sslmode=disable"
	dbTest, err := sqlx.Open("pgx", dsn)
	require.Nil(t, err)

	err = dbTest.PingContext(ctx)
	require.Nil(t, err)

	dbTest.MustExec(schema)

	dbTest.Close()

	st := New()

	err = st.Connect(ctx, dsn)
	require.Nil(t, err)
	defer func() {
		if err := st.Close(ctx); err != nil {
			fmt.Printf("cannot close psql connection: " + err.Error())
		}
	}()

	err = st.Add(storage.Event{
		ID:        "1",
		Title:     "event 1",
		TimeStart: time.Date(2022, time.Month(1), 1, 10, 12, 9, 0, time.UTC),
	})
	require.Nil(t, err)

	err = st.Add(storage.Event{
		ID:               "2",
		Title:            "event 2",
		TimeStart:        time.Date(2022, time.Month(1), 1, 12, 10, 0, 0, time.UTC),
		NotificationTime: 5 * time.Minute,
	})
	require.Nil(t, err)

	err = st.Add(storage.Event{
		ID:               "3",
		Title:            "event 3",
		TimeStart:        time.Date(2022, time.Month(1), 1, 13, 10, 0, 0, time.UTC),
		NotificationTime: 5 * time.Minute,
	})

	require.Nil(t, err)

	res, err := st.GetEventsForNotify(time.Date(2022, time.Month(1), 1, 12, 10, 0, 0, time.UTC))
	require.Nil(t, err)
	require.Equal(t, 1, len(res))

	require.Equal(t, "2", res[0].ID)

	err = st.MarkEventIsNotifyed(res[0].ID)
	require.Nil(t, err)

	res, err = st.GetEventsForNotify(time.Date(2022, time.Month(1), 1, 12, 10, 0, 0, time.UTC))
	require.Nil(t, err)
	require.Equal(t, 0, len(res))

}
