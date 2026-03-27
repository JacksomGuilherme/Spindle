package utils

import (
	"sync"
	"time"
)

type Pairing struct {
	Authenticated bool
	UserID        string
	ExpiresAt     time.Time
}

type PairingStore struct {
	mu       sync.RWMutex
	pairings map[string]*Pairing
}

func NewPairingStore() *PairingStore {
	return &PairingStore{
		pairings: make(map[string]*Pairing),
	}
}

func (s *PairingStore) Create(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pairings[id] = &Pairing{
		Authenticated: false,
		ExpiresAt:     time.Now().Add(5 * time.Minute),
	}
}

func (s *PairingStore) Get(id string) (*Pairing, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.pairings[id]
	return p, ok
}

func (s *PairingStore) Authenticate(id string, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p, ok := s.pairings[id]; ok {
		p.Authenticated = true
		p.UserID = userID
	}
}

func (s *PairingStore) Cleanup() {
	for {
		time.Sleep(1 * time.Minute)

		s.mu.Lock()
		for id, p := range s.pairings {
			if time.Now().After(p.ExpiresAt) {
				delete(s.pairings, id)
			}
		}
		s.mu.Unlock()
	}
}
