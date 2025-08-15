package internal

import "sync"

type Store struct {
    mu   sync.RWMutex
    data map[string]string
}

func NewStore() *Store {
    return &Store{data: make(map[string]string)}
}

func (s *Store) Set(key, value string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    v, ok := s.data[key]
    return v, ok
}
