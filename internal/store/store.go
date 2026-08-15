package store

import (
	"errors"
	"sync"

	"mailer/internal/model"
)

var (
	ErrNotFound      = errors.New("message not found")
	ErrEmpty         = errors.New("queue empty")
	ErrAlreadyExists = errors.New("message already exists")
)

type Store struct {
	mu       sync.RWMutex
	messages map[string]*model.Message
	order    []string
	dedup    map[string]bool
}

func New() *Store {
	return &Store{
		messages: make(map[string]*model.Message),
		order:    []string{},
		dedup:    make(map[string]bool),
	}
}

func (s *Store) Enqueue(m *model.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.messages[m.ID]; ok {
		return ErrAlreadyExists
	}
	s.messages[m.ID] = m
	s.order = append(s.order, m.ID)
	s.dedup[m.ID] = true
	return nil
}

func (s *Store) Get(id string) (*model.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.messages[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

// Pending 返回所有 queued 消息，按优先级排好序（返回全新切片）。
func (s *Store) Pending() []*model.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Message, 0, len(s.order))
	for _, id := range s.order {
		m := s.messages[id]
		if m.Status == model.StatusQueued {
			out = append(out, m)
		}
	}
	return out
}

func (s *Store) PendingIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.order))
	copy(out, s.order)
	return out
}

func (s *Store) MarkSent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.messages[id]
	if !ok {
		return ErrNotFound
	}
	m.Status = model.StatusSent
	return nil
}

func (s *Store) MarkFailed(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.messages[id]
	if !ok {
		return ErrNotFound
	}
	m.Attempts++
	m.Status = model.StatusFailed
	return nil
}
