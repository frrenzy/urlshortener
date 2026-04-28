package repository

import (
	"sync"

	"frrenzy/urlshortener/internal/model"
)

type mapStorage struct {
	Repository
	storage map[string]model.Link
	mu      *sync.RWMutex
}

func (s mapStorage) Add(l model.Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[l.ShortURL] = l

	return nil
}

func (s mapStorage) Get(short string) (*model.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	link, exists := s.storage[short]
	if !exists {
		return &model.Link{}, errNotFound
	}

	return &link, nil
}

func NewMapStorage() mapStorage {
	instance := make(map[string]model.Link)

	return mapStorage{
		storage: instance,
		mu:      &sync.RWMutex{},
	}
}
