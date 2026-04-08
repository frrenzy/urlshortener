// Package repository
package repository

import (
	"errors"
	"net/url"
)

type mapStorage struct {
	Repository
	storage map[string]url.URL
}

func (s mapStorage) Add(original url.URL, short string) {
	s.storage[short] = original
}

func (s mapStorage) Get(short string) (url.URL, error) {
	original, exists := s.storage[short]
	if !exists {
		return url.URL{}, errors.New("not found")
	}

	return original, nil
}

func NewMapStorage() mapStorage {
	instance := make(map[string]url.URL)

	return mapStorage{
		storage: instance,
	}
}
