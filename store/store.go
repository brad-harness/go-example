package store

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("paste not found")
)

type Paste struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

type Store interface {
	Create(content string, ttl time.Duration) (*Paste, error)
	Get(id string) (*Paste, error)
	Delete(id string) error
	List() ([]*Paste, error)
}

type MemoryStore struct {
	mu     sync.RWMutex
	pastes map[string]*Paste
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		pastes: make(map[string]*Paste),
	}
}

func (s *MemoryStore) Create(content string, ttl time.Duration) (*Paste, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	paste := &Paste{
		ID:        uuid.New().String(),
		Content:   content,
		CreatedAt: time.Now(),
	}

	if ttl > 0 {
		paste.ExpiresAt = time.Now().Add(ttl)
	}

	s.pastes[paste.ID] = paste
	return paste, nil
}

func (s *MemoryStore) Get(id string) (*Paste, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	paste, ok := s.pastes[id]
	if !ok {
		return nil, ErrNotFound
	}

	if !paste.ExpiresAt.IsZero() && time.Now().After(paste.ExpiresAt) {
		return nil, ErrNotFound
	}

	return paste, nil
}

func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.pastes[id]; !ok {
		return ErrNotFound
	}

	delete(s.pastes, id)
	return nil
}

func (s *MemoryStore) List() ([]*Paste, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pastes := make([]*Paste, 0, len(s.pastes))
	for _, paste := range s.pastes {
		if paste.ExpiresAt.IsZero() || time.Now().Before(paste.ExpiresAt) {
			pastes = append(pastes, paste)
		}
	}

	return pastes, nil
}
