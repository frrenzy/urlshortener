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

type URLBatchInstance struct {
	Original      url.URL
	CorrelationID string
}

func (s ShortenerService) CreateShortURL(ctx context.Context, original url.URL) (string, error) {
	short := s.generator.generateShort()

	err := s.storage.Add(ctx, model.NewLink(original, short))
	if err != nil {
		return "", err
	}

	return short, nil
}

func (s ShortenerService) GetOriginalURL(ctx context.Context, short string) (string, error) {
	link, err := s.storage.Get(ctx, short)
	if err != nil {
		return "", err
	}

	return link.OriginalURL.String(), nil
}

func (s ShortenerService) BatchCreateShortURL(ctx context.Context, urls []URLBatchInstance) ([]model.Link, error) {
	var links []model.Link
	for _, instance := range urls {
		short := s.generator.generateShort()
		links = append(links, model.NewLinkWithUUID(instance.Original, short, instance.CorrelationID))
	}

	result, err := s.storage.BatchAdd(ctx, links)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s ShortenerService) PingStorage(ctx context.Context) error {
	return s.storage.PingStorage(ctx)
}

func NewShortenerService(repository repository.Repository) ShortenerService {
	return ShortenerService{
		storage:   repository,
		generator: randomGenerator{},
	}
}
