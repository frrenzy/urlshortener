package repository

import (
	"context"
	"sync"

	"frrenzy/urlshortener/internal/model"
)

type mapStorage struct {
	Repository
	storage map[string]model.Link
	mu      *sync.RWMutex
}

func (s mapStorage) Add(ctx context.Context, l model.Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[l.ShortURL] = l

	return nil
}

func (s mapStorage) BatchAdd(ctx context.Context, links []model.Link) ([]model.Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, link := range links {
		s.storage[link.ShortURL] = link
	}

	return links, nil
}

func (s mapStorage) Get(ctx context.Context, short string) (*model.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	link, exists := s.storage[short]
	if !exists {
		return &model.Link{}, errNotFound
	}

	return &link, nil
}

func (s mapStorage) Close() {}

func (s mapStorage) PingStorage(ctx context.Context) error {
	return errDBNotConnected
}

func NewMapStorage() mapStorage {
	instance := make(map[string]model.Link)

	return mapStorage{
		storage: instance,
		mu:      &sync.RWMutex{},
	}
}
