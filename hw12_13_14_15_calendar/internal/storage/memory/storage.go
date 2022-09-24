package memorystorage

import (
	"sync"

	"github.com/julinserg/go_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	// TODO
	mu     sync.RWMutex
	events map[string]storage.Event
}

func (s *Storage) Get(id string) storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events[id]
}

func (s *Storage) Add(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(event.ID) == 0 {
		return storage.ErrEventIdNotSet
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) Update(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(event.ID) == 0 {
		return storage.ErrEventIdNotSet
	}
	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrEventIdNotExist
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) Remove(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; !ok {
		return storage.ErrEventIdNotExist
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) create() error {
	s.events = make(map[string]storage.Event)
	return nil
}

func New() *Storage {
	s := &Storage{}
	s.create()
	return s
}

// TODO
