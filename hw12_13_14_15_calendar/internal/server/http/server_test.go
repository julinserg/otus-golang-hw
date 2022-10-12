package internalhttp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
	}{
		{"bad_request", http.MethodPost, "http://test.test", nil, http.StatusBadRequest},
		{
			"add-event-1",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "1", "title": "testTitle", "description": "testDescription", "time_start": "2021-02-18T21:54:42.123Z"}`),
			http.StatusOK,
		},
		{
			"add-event-2",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "2", "title": "test", "description": "testDescription", "time_start": "2022-02-18T21:54:42.123Z"}`),
			http.StatusOK,
		},
		{
			"add-event-3",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "3", "title": "test", "description": "testDescription", "time_start": "2022-02-18T21:54:42.123Z"}`),
			http.StatusInternalServerError,
		},
		{
			"add-event-4",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "4", "title": "test}`),
			http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(c.method, c.target, c.body)
			w := httptest.NewRecorder()
			service.addEvent(w, r)
			resp := w.Result()
			require.Equal(t, c.responseCode, resp.StatusCode)
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
	}{
		{"bad_request", http.MethodPost, "http://test.test", nil, http.StatusBadRequest},
		{
			"remove-event-1",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "1"}`),
			http.StatusOK,
		},
		{
			"remove-event-2",
			http.MethodPost,
			"http://test.test",
			bytes.NewBufferString(`{"id": "1"}`),
			http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(c.method, c.target, c.body)
			w := httptest.NewRecorder()
			service.removeEvent(w, r)
			resp := w.Result()
			require.Equal(t, c.responseCode, resp.StatusCode)
		})
	}
}
