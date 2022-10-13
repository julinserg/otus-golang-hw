package internalhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
	stor "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

type LoggerFakeImpl struct {
}

func (l LoggerFakeImpl) Info(msg string) {
	fmt.Println("Info: " + msg)
}

func (l LoggerFakeImpl) Error(msg string) {
	fmt.Println("Error: " + msg)
}

func (l LoggerFakeImpl) Debug(msg string) {
	fmt.Println("Debug: " + msg)
}

func (l LoggerFakeImpl) Warn(msg string) {
	fmt.Println("Warn: " + msg)
}

func TestServiceAddEvent(t *testing.T) {
	logg := LoggerFakeImpl{}
	storage := memorystorage.New()
	calendar := app.New(logg, storage)
	service := calendarHandler{logg, calendar}

	cases := []struct {
		name         string
		method       string
		target       string
		body         io.Reader
		responseCode int
		responseBody Response
	}{
		{"bad_request", http.MethodPost, "http://test.test", nil, http.StatusBadRequest,
			Response{Error: struct {
				Message string `json:"message"`
			}{
				Message: "unexpected end of JSON input",
			}},
		},
		{
			"add-event-1",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "1", "title": "testTitle", "description": "testDescription", "time_start": "2021-02-18T21:54:42.123Z"}`),
			http.StatusOK,
			Response{},
		},
		{
			"add-event-2",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "2", "title": "test", "description": "testDescription", "time_start": "2022-02-18T21:54:42.123Z"}`),
			http.StatusOK,
			Response{},
		},
		{
			"add-event-3",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "3", "title": "test", "description": "testDescription", "time_start": "2022-02-18T21:54:42.123Z"}`),
			http.StatusInternalServerError,
			Response{Error: struct {
				Message string `json:"message"`
			}{
				Message: "time event is busy",
			}},
		},
		{
			"add-event-4",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "4", "title": "test}`),
			http.StatusBadRequest,
			Response{Error: struct {
				Message string `json:"message"`
			}{
				Message: "unexpected end of JSON input",
			}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(c.method, c.target, c.body)
			w := httptest.NewRecorder()
			service.addEvent(w, r)
			resp := w.Result()
			require.Equal(t, c.responseCode, resp.StatusCode)

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)
			result := Response{}
			json.Unmarshal(body, &result)
			require.Equal(t, result, c.responseBody)
		})
	}
}

func TestServiceRemoveEvent(t *testing.T) {
	logg := LoggerFakeImpl{}
	storage := memorystorage.New()
	calendar := app.New(logg, storage)
	service := calendarHandler{logg, calendar}

	event1 := stor.Event{
		ID: "1", Title: "event 1",
		TimeStart: time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC),
	}
	err := storage.Add(event1)
	require.Nil(t, err)

	cases := []struct {
		name         string
		method       string
		target       string
		body         io.Reader
		responseCode int
		responseBody Response
	}{
		{"bad_request", http.MethodPost, "http://test.test", nil, http.StatusBadRequest,
			Response{Error: struct {
				Message string `json:"message"`
			}{
				Message: "unexpected end of JSON input",
			}}},
		{
			"remove-event-1",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "1"}`),
			http.StatusOK,
			Response{},
		},
		{
			"remove-event-2",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "1"}`),
			http.StatusInternalServerError,
			Response{Error: struct {
				Message string `json:"message"`
			}{
				Message: "Event ID not exist",
			}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(c.method, c.target, c.body)
			w := httptest.NewRecorder()
			service.removeEvent(w, r)
			resp := w.Result()
			require.Equal(t, c.responseCode, resp.StatusCode)

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)
			result := Response{}
			json.Unmarshal(body, &result)
			require.Equal(t, result, c.responseBody)
		})
	}
}

func TestServiceUpdateEvent(t *testing.T) {
	logg := LoggerFakeImpl{}
	storage := memorystorage.New()
	calendar := app.New(logg, storage)
	service := calendarHandler{logg, calendar}

	event1 := stor.Event{
		ID: "1", Title: "event 1",
		TimeStart: time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC),
	}
	err := storage.Add(event1)
	require.Nil(t, err)

	cases := []struct {
		name         string
		method       string
		target       string
		body         io.Reader
		responseCode int
		responseBody Response
	}{
		{"bad_request", http.MethodPost, "http://test.test", nil, http.StatusBadRequest,
			Response{Error: struct {
				Message string `json:"message"`
			}{
				Message: "unexpected end of JSON input",
			}}},
		{
			"update-event-1",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "1", "title": "testTitle", "description": "testDescription", "time_start": "2021-02-18T21:54:42.123Z"}`),
			http.StatusOK,
			Response{},
		},
		{
			"update-event-2",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "2", "title": "testTitle", "description": "testDescription", "time_start": "2021-02-18T21:54:42.123Z"}`),
			http.StatusInternalServerError,
			Response{Error: struct {
				Message string `json:"message"`
			}{
				Message: "Event ID not exist",
			}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(c.method, c.target, c.body)
			w := httptest.NewRecorder()
			service.updateEvent(w, r)
			resp := w.Result()
			require.Equal(t, c.responseCode, resp.StatusCode)

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)
			result := Response{}
			json.Unmarshal(body, &result)
			require.Equal(t, result, c.responseBody)
		})
	}
}

func TestServiceGetEvents(t *testing.T) {
	logg := LoggerFakeImpl{}
	storage := memorystorage.New()
	calendar := app.New(logg, storage)
	service := calendarHandler{logg, calendar}

	event1 := stor.Event{
		ID: "1", Title: "event 1",
		TimeStart: time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC),
	}
	tt := time.Date(2022, time.Month(1), 1, 1, 10, 30, 0, time.UTC)
	str := tt.String()
	_ = str
	err := storage.Add(event1)
	require.Nil(t, err)

	event2 := stor.Event{
		ID: "2", Title: "event 2",
		TimeStart: time.Date(2022, time.Month(1), 3, 1, 10, 30, 0, time.UTC),
	}
	err = storage.Add(event2)
	require.Nil(t, err)

	event3 := stor.Event{
		ID: "3", Title: "event 3",
		TimeStart: time.Date(2022, time.Month(1), 14, 1, 10, 30, 0, time.UTC),
	}
	err = storage.Add(event3)
	require.Nil(t, err)

	cases := []struct {
		name         string
		method       string
		target       string
		body         io.Reader
		responseCode int
		responseBody Response
	}{
		{"bad_request", http.MethodGet, "http://test.test/get_by_day", nil, http.StatusBadRequest,
			Response{Error: struct {
				Message string `json:"message"`
			}{
				Message: "unexpected end of JSON input",
			}}},
		{
			"get-event-by-day-1",
			http.MethodGet,
			"http://test.test/get_by_day",
			bytes.NewBufferString(`{"time": "2022-01-01T00:00:00Z"}`),
			http.StatusOK,
			Response{Data: []app.EventApp{{ID: "1",
				Title:            "event 1",
				TimeStart:        time.Date(2022, time.January, 1, 1, 10, 30, 0, time.UTC),
				TimeEnd:          time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
				Description:      "",
				UserID:           "",
				NotificationTime: 0}},
			},
		},
		{
			"get-event-by-day-2",
			http.MethodGet,
			"http://test.test/get_by_day",
			bytes.NewBufferString(`{"time": "2023-01-01T00:00:00Z"}`),
			http.StatusOK,
			Response{Data: []app.EventApp{}},
		},
		{
			"get-event-by-week-1",
			http.MethodGet,
			"http://test.test/get_by_week",
			bytes.NewBufferString(`{"time": "2022-01-01T00:00:00Z"}`),
			http.StatusOK,
			Response{Data: []app.EventApp{
				{
					ID:               "1",
					Title:            "event 1",
					TimeStart:        time.Date(2022, time.January, 1, 1, 10, 30, 0, time.UTC),
					TimeEnd:          time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					Description:      "",
					UserID:           "",
					NotificationTime: 0,
				},
				{
					ID:               "2",
					Title:            "event 2",
					TimeStart:        time.Date(2022, time.January, 3, 1, 10, 30, 0, time.UTC),
					TimeEnd:          time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					Description:      "",
					UserID:           "",
					NotificationTime: 0,
				}}},
		},
		{
			"get-event-by-week-2",
			http.MethodGet,
			"http://test.test/get_by_week",
			bytes.NewBufferString(`{"time": "2023-01-01T00:00:00Z"}`),
			http.StatusOK,
			Response{Data: []app.EventApp{}},
		},
		{
			"get-event-by-month-1",
			http.MethodGet,
			"http://test.test/get_by_month",
			bytes.NewBufferString(`{"time": "2022-01-01T00:00:00Z"}`),
			http.StatusOK,
			Response{Data: []app.EventApp{
				{
					ID:               "1",
					Title:            "event 1",
					TimeStart:        time.Date(2022, time.January, 1, 1, 10, 30, 0, time.UTC),
					TimeEnd:          time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					Description:      "",
					UserID:           "",
					NotificationTime: 0,
				},
				{
					ID:               "2",
					Title:            "event 2",
					TimeStart:        time.Date(2022, time.January, 3, 1, 10, 30, 0, time.UTC),
					TimeEnd:          time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					Description:      "",
					UserID:           "",
					NotificationTime: 0,
				},
				{
					ID:               "3",
					Title:            "event 3",
					TimeStart:        time.Date(2022, time.January, 14, 1, 10, 30, 0, time.UTC),
					TimeEnd:          time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
					Description:      "",
					UserID:           "",
					NotificationTime: 0,
				}}},
		},
		{
			"get-event-by-month-2",
			http.MethodGet,
			"http://test.test/get_by_month",
			bytes.NewBufferString(`{"time": "2023-01-01T00:00:00Z"}`),
			http.StatusOK,
			Response{Data: []app.EventApp{}},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(c.method, c.target, c.body)
			w := httptest.NewRecorder()
			if c.target == "http://test.test/get_by_day" {
				service.getEventsByDay(w, r)
			} else if c.target == "http://test.test/get_by_week" {
				service.getEventsByWeek(w, r)
			} else if c.target == "http://test.test/get_by_month" {
				service.getEventsByMonth(w, r)
			}
			resp := w.Result()
			require.Equal(t, c.responseCode, resp.StatusCode)

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)
			result := Response{}
			json.Unmarshal(body, &result)

			if !reflect.DeepEqual(c.responseBody, result) {
				t.Fatalf("[%s] results not match\nGot : %#v\nWant: %#v", c.name, result, c.responseBody)
				return
			}
		})
	}
}
