package memorystorage

import (
	"strconv"
	"sync"
	"testing"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {

	t.Run("basic", func(t *testing.T) {
		st := New()

		event1 := storage.Event{ID: "{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}", Title: "event 1"}
		st.Add(event1)
		event2 := storage.Event{ID: "{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}", Title: "event 2"}
		st.Add(event2)

		require.Equal(t, event1, st.Get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}"))
		require.Equal(t, event2, st.Get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}"))
		require.Equal(t, storage.Event{}, st.Get("{81e125ce-072e-4556-8a4c-597572a7277a}"))

		require.Equal(t, "event 1", st.Get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}").Title)
		require.Equal(t, "event 2", st.Get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}").Title)

		event1.Title = "event 5"
		st.Update(event1)

		require.Equal(t, "event 5", st.Get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}").Title)
		require.Equal(t, "event 2", st.Get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}").Title)

		st.Remove(event1.ID)
		require.Equal(t, storage.Event{}, st.Get("{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}"))
		require.Equal(t, event2, st.Get("{95c3d43f-a8be-49ee-b5c6-d98fb25a38bc}"))

		event1.ID = "{0e745e54-0f24-4b4f-aa9f-3bd1167e55f9}"
		event1.Title = "event 6"
		require.ErrorIs(t, st.Update(event1), storage.ErrEventIdNotExist)

		require.ErrorIs(t, st.Remove(event1.ID), storage.ErrEventIdNotExist)

		event3 := storage.Event{Title: "event 3"}
		require.ErrorIs(t, st.Add(event3), storage.ErrEventIdNotSet)
	})
	t.Run("goroutine", func(t *testing.T) {
		st := New()
		wg := &sync.WaitGroup{}
		wg.Add(20)
		for i := 0; i < 10; i++ {
			go func(i int) {
				defer wg.Done()
				event := storage.Event{ID: strconv.Itoa(i), Title: strconv.Itoa(i)}
				st.Add(event)
			}(i)
			go func(i int) {
				defer wg.Done()
				st.Get(strconv.Itoa(i))
			}(i)
		}
		wg.Wait()
		require.Equal(t, "0", st.Get("0").Title)
		require.Equal(t, "1", st.Get("1").Title)
		require.Equal(t, "2", st.Get("2").Title)
		require.Equal(t, "3", st.Get("3").Title)
		require.Equal(t, "4", st.Get("4").Title)
		require.Equal(t, "5", st.Get("5").Title)
		require.Equal(t, "6", st.Get("6").Title)
		require.Equal(t, "7", st.Get("7").Title)
		require.Equal(t, "8", st.Get("8").Title)
		require.Equal(t, "9", st.Get("9").Title)
	})

}
