package memorystorage

import (
	"sync"
	"time"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	// TODO
	mu         sync.RWMutex
	events     map[string]storage.Event
	eventsTime map[time.Time]string
}

func (s *Storage) get(id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var ev storage.Event
	var ok bool
	if ev, ok = s.events[id]; !ok {
		return ev, storage.ErrEventIDNotExist
	}
	return ev, nil
}

func dateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return dateEqual(check, start)
	}
	return !start.After(check) || !end.Before(check)
}

func (s *Storage) getEventsBy(date1, date2 time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storage.Event, 0)
	for _, v := range s.events {
		if inTimeSpan(date1, date2, v.TimeStart) {
			result = append(result, v)
		}
	}
	return result, nil
}

func (s *Storage) GetEventsByDay(date time.Time) ([]storage.Event, error) {
	return s.getEventsBy(date, date)
}

func (s *Storage) GetEventsByWeek(dateBeginWeek time.Time) ([]storage.Event, error) {
	dateEndWeek := dateBeginWeek.AddDate(0, 0, 7)
	return s.getEventsBy(dateBeginWeek, dateEndWeek)
}

func (s *Storage) GetEventsByMonth(dateBeginMonth time.Time) ([]storage.Event, error) {
	dateEndMonth := dateBeginMonth.AddDate(0, 1, 0)
	return s.getEventsBy(dateBeginMonth, dateEndMonth)
}

func (s *Storage) Add(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(event.ID) == 0 {
		return storage.ErrEventIDNotSet
	}
	if _, ok := s.events[event.ID]; ok {
		return storage.ErrEventIDAlreadyExist
	}
	if _, ok := s.eventsTime[event.TimeStart]; ok {
		return storage.ErrTimeBusy
	}
	s.events[event.ID] = event
	s.eventsTime[event.TimeStart] = event.ID
	return nil
}

func (s *Storage) Update(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(event.ID) == 0 {
		return storage.ErrEventIDNotSet
	}
	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrEventIDNotExist
	}
	if _, ok := s.eventsTime[event.TimeStart]; ok && s.events[event.ID].TimeStart != event.TimeStart {
		return storage.ErrTimeBusy
	}
	s.events[event.ID] = event
	s.eventsTime[event.TimeStart] = event.ID
	return nil
}

func (s *Storage) Remove(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ev, ok := s.events[id]
	if !ok {
		return storage.ErrEventIDNotExist
	}
	delete(s.eventsTime, ev.TimeStart)
	delete(s.events, id)
	return nil
}

func (s *Storage) create() {
	s.events = make(map[string]storage.Event)
	s.eventsTime = make(map[time.Time]string)
}

func New() *Storage {
	s := &Storage{}
	s.create()
	return s
}
