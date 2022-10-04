package memorystorage

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorageBasic(t *testing.T) {
	st := New()

	event1 := storage.Event{
		ID: "{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}", Title: "event 1",
		TimeStart: time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC),
	}
	err := st.Add(event1)
	require.Nil(t, err)
	event2 := storage.Event{
		ID: "{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}", Title: "event 2",
		TimeStart: time.Date(2022, time.Month(2), 1, 1, 10, 30, 0, time.UTC),
	}
	err = st.Add(event2)
	require.Nil(t, err)

	resGet, _ := st.get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}")
	require.Equal(t, event1, resGet)
	resGet, _ = st.get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}")
	require.Equal(t, event2, resGet)
	resGet, _ = st.get("{81e125ce-072e-4556-8a4c-597572a7277a}")
	require.Equal(t, storage.Event{}, resGet)

	resGet, _ = st.get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}")
	require.Equal(t, "event 1", resGet.Title)
	resGet, _ = st.get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}")
	require.Equal(t, "event 2", resGet.Title)

	event1.Title = "event 5"
	st.Update(event1)

	resGet, _ = st.get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}")
	require.Equal(t, "event 5", resGet.Title)
	resGet, _ = st.get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}")
	require.Equal(t, "event 2", resGet.Title)

	st.Remove(event1.ID)
	resGet, _ = st.get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}")
	require.Equal(t, storage.Event{}, resGet)
	resGet, _ = st.get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}")
	require.Equal(t, event2, resGet)

	event1.ID = "{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}"
	event1.Title = "event 6"
	require.ErrorIs(t, st.Update(event1), storage.ErrEventIDNotExist)

	require.ErrorIs(t, st.Remove(event1.ID), storage.ErrEventIDNotExist)

	event3 := storage.Event{Title: "event 3"}
	require.ErrorIs(t, st.Add(event3), storage.ErrEventIDNotSet)
}

func TestStorageGoroutine(t *testing.T) {
	st := New()
	wg := &sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			event := storage.Event{
				ID:        strconv.Itoa(i),
				Title:     strconv.Itoa(i),
				TimeStart: time.Date(2022, time.Month(1), i+1, 1, 10, 30, 0, time.UTC),
			}
			st.Add(event)
		}(i)
		go func(i int) {
			defer wg.Done()
			st.get(strconv.Itoa(i))
		}(i)
	}
	wg.Wait()
	resGet, _ := st.get("0")
	require.Equal(t, "0", resGet.Title)
	resGet, _ = st.get("1")
	require.Equal(t, "1", resGet.Title)
	resGet, _ = st.get("2")
	require.Equal(t, "2", resGet.Title)
	resGet, _ = st.get("3")
	require.Equal(t, "3", resGet.Title)
	resGet, _ = st.get("4")
	require.Equal(t, "4", resGet.Title)
	resGet, _ = st.get("5")
	require.Equal(t, "5", resGet.Title)
	resGet, _ = st.get("6")
	require.Equal(t, "6", resGet.Title)
	resGet, _ = st.get("7")
	require.Equal(t, "7", resGet.Title)
	resGet, _ = st.get("8")
	require.Equal(t, "8", resGet.Title)
	resGet, _ = st.get("9")
	require.Equal(t, "9", resGet.Title)
}

func TestStorageLogic(t *testing.T) {
	st := New()

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
