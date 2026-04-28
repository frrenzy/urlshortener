// Package service
package service

import (
	"net/url"

	"frrenzy/urlshortener/internal/model"
	"frrenzy/urlshortener/internal/repository"
)

type ShortenerService struct {
	storage   repository.Repository
	generator shortener
}

func (s ShortenerService) CreateShortURL(original url.URL) (string, error) {
	short := s.generator.generateShort()

	s.storage.Add(model.NewLink(original, short))

	return short, nil
}

func (s ShortenerService) GetOriginalURL(short string) (string, error) {
	link, err := s.storage.Get(short)
	if err != nil {
		return "", err
	}

	return link.OriginalURL.String(), nil
}

func NewShortenerService(repository repository.Repository) ShortenerService {
	return ShortenerService{
		storage:   repository,
		generator: randomGenerator{},
	}
}
