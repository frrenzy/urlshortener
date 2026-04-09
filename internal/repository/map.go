package repository

import (
	"sync"

	"frrenzy/urlshortener/internal/model"
)

type mapStorage struct {
	Repository
	storage map[string]model.Link
	lock    *sync.RWMutex
}

func (s mapStorage) Add(l model.Link) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.storage[l.ShortURL] = l

	return nil
}

func (s mapStorage) Get(short string) (model.Link, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	link, exists := s.storage[short]
	if !exists {
		return model.Link{}, errNotFound
	}

	return link, nil
}

func NewMapStorage() mapStorage {
	instance := make(map[string]model.Link)

	return mapStorage{
		storage: instance,
		lock:    &sync.RWMutex{},
	}
}
