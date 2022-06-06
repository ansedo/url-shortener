package memorystorage

import (
	"errors"
	"sync"
)

type Storage struct {
	sync.RWMutex
	repo map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		repo: make(map[string]string),
	}
}

func (s *Storage) Get(key string) (string, error) {
	if !s.Has(key) {
		return "", errors.New("this key does not exist")
	}
	s.RLock()
	defer s.RUnlock()
	return s.repo[key], nil
}

func (s *Storage) Set(key, value string) error {
	if s.Has(key) {
		return errors.New("this key already exists")
	}
	s.Lock()
	defer s.Unlock()
	s.repo[key] = value
	return nil
}

func (s *Storage) Has(key string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.repo[key]
	return ok
}

func (s *Storage) NextID() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.repo)
}
