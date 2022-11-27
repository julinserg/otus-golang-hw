package internalhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/app"
)

func hellowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is my calendar!"))
}

type RemoveEvent struct {
	ID string `json:"id"`
}

type GetEvent struct {
	Time time.Time `json:"time"`
}

type calendarHandler struct {
	logger Logger
	app    Application
}

type Response struct {
	Data  []app.Event `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (ch *calendarHandler) WriteResponse(w http.ResponseWriter, resp *Response) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		ch.logger.Error("response marshal error: " + err.Error())
	}
	_, err = w.Write(resBuf)
	if err != nil {
		ch.logger.Error("response marshal error: " + err.Error())
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return
}

func (ch *calendarHandler) checkRequest(methodType string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != methodType {
		resp := &Response{}
		resp.Error.Message = fmt.Sprintf("method %s not supported on uri %s", r.Method, r.URL.Path)
		ch.logger.Error(resp.Error.Message)
		w.WriteHeader(http.StatusMethodNotAllowed)
		ch.WriteResponse(w, resp)
		return false
	}
	return true
}

func (ch *calendarHandler) checkErrorAndSendResponse(err error, code int, w http.ResponseWriter) bool {
	if err != nil {
		resp := &Response{}
		resp.Error.Message = err.Error()
		ch.logger.Error(resp.Error.Message)
		w.WriteHeader(code)
		ch.WriteResponse(w, resp)
		return false
	}
	return true
}

type (
	actionPost func() error
	actionGet  func(date time.Time) ([]app.Event, error)
)

func (ch *calendarHandler) genericHandlerPost(w http.ResponseWriter, r *http.Request, data interface{}, act actionPost) bool {

	if !ch.checkRequest(http.MethodPost, w, r) {
		return false
	}

	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		resp := &Response{}
		resp.Error.Message = err.Error()
		ch.logger.Error(resp.Error.Message)
		w.WriteHeader(http.StatusBadRequest)
		ch.WriteResponse(w, resp)
		return false
	}

	err = json.Unmarshal(buf, data)
	if !ch.checkErrorAndSendResponse(err, http.StatusBadRequest, w) {
		return false
	}
	err = act()
	if !ch.checkErrorAndSendResponse(err, http.StatusInternalServerError, w) {
		return false
	}
	w.WriteHeader(http.StatusOK)
	return true
}

func (ch *calendarHandler) genericHandlerGet(w http.ResponseWriter, r *http.Request, act actionGet) bool {

	if !ch.checkRequest(http.MethodGet, w, r) {
		return false
	}
	strTime := r.URL.Query().Get("time")
	var err error
	if len(strTime) == 0 {
		err = fmt.Errorf("Time is empty")
	}

	time, err := time.Parse(time.RFC3339, strTime)
	if !ch.checkErrorAndSendResponse(err, http.StatusBadRequest, w) {
		return false
	}
	events, err := act(time)
	if !ch.checkErrorAndSendResponse(err, http.StatusInternalServerError, w) {
		return false
	}
	w.WriteHeader(http.StatusOK)
	resp := &Response{Data: events}
	ch.WriteResponse(w, resp)
	return true
}

func (ch *calendarHandler) addEvent(w http.ResponseWriter, r *http.Request) {
	req := &app.Event{}
	ch.genericHandlerPost(w, r, req, func() error { return ch.app.AddEvent(req) })
}

func (ch *calendarHandler) removeEvent(w http.ResponseWriter, r *http.Request) {
	req := &RemoveEvent{}
	ch.genericHandlerPost(w, r, req, func() error { return ch.app.RemoveEvent(req.ID) })
}

func (ch *calendarHandler) updateEvent(w http.ResponseWriter, r *http.Request) {
	req := &app.Event{}
	ch.genericHandlerPost(w, r, req, func() error { return ch.app.UpdateEvent(req) })
}

func (ch *calendarHandler) getEventsByDay(w http.ResponseWriter, r *http.Request) {
	ch.genericHandlerGet(w, r, func(date time.Time) ([]app.Event, error) { return ch.app.GetEventsByDay(date) })
}

func (ch *calendarHandler) getEventsByMonth(w http.ResponseWriter, r *http.Request) {
	ch.genericHandlerGet(w, r, func(date time.Time) ([]app.Event, error) { return ch.app.GetEventsByMonth(date) })
}

func (ch *calendarHandler) getEventsByWeek(w http.ResponseWriter, r *http.Request) {
	ch.genericHandlerGet(w, r, func(date time.Time) ([]app.Event, error) { return ch.app.GetEventsByWeek(date) })
}
