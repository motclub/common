package helpers

import "sync"

type safeMapString struct {
	m map[string]string
	l *sync.RWMutex
}

func (s *safeMapString) Set(key string, value string) {
	s.l.Lock()
	defer s.l.Unlock()
	s.m[key] = value
}

func (s *safeMapString) Get(key string) string {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.m[key]
}

func NewSafeMapString() *safeMapString {
	return &safeMapString{l: new(sync.RWMutex), m: make(map[string]string)}
}
