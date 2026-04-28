// Package service
package service

import (
	"context"
	"net/url"

	"frrenzy/urlshortener/internal/model"
	"frrenzy/urlshortener/internal/repository"
)

type ShortenerService struct {
	storage   repository.Repository
	generator shortener
}

func (s ShortenerService) CreateShortURL(ctx context.Context, original url.URL) (string, error) {
	short := s.generator.generateShort()

	s.storage.Add(ctx, model.NewLink(original, short))

	return short, nil
}

func (s ShortenerService) GetOriginalURL(ctx context.Context, short string) (string, error) {
	link, err := s.storage.Get(ctx, short)
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
