package auth

import "sync"

type TokenStore interface {
	Current() int64
	Next() int64
	Set(value int64)
	Reset()
}

type MemoryTokenStore struct {
	mu      sync.Mutex
	current int64
}

func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{}
}

func (s *MemoryTokenStore) Current() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}

func (s *MemoryTokenStore) Next() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current++
	return s.current
}

func (s *MemoryTokenStore) Set(value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = value
}

func (s *MemoryTokenStore) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.current = 0
}
