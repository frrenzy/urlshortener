// Package service
package service

import (
	"context"
	"errors"
	"net/url"
	"time"

	"frrenzy/urlshortener/internal/model"
	"frrenzy/urlshortener/internal/repository"
	"frrenzy/urlshortener/internal/util/logger"

	"go.uber.org/zap"
)

type deleteTask struct {
	userID    int
	shortURLs []string
}

type ShortenerService struct {
	storage           repository.Repository
	generator         shortener
	deleteJobsChannel chan deleteTask
	deleteTicker      *time.Ticker
}

type URLBatchInstance struct {
	Original      url.URL
	CorrelationID string
}

func (s ShortenerService) CreateShortURL(ctx context.Context, original url.URL, userID int) (string, error) {
	short := s.generator.generateShort()

	err := s.storage.Add(ctx, model.NewLink(original, short, userID))
	if err != nil {
		if !errors.Is(err, repository.ErrDBExisting) {
			return "", err
		}

		link, err := s.storage.GetByOriginal(ctx, original)
		if err != nil {
			return "", err
		}

		return link.ShortURL, ErrLinkAlreadyExists
	}

	return short, nil
}

func (s ShortenerService) GetOriginalURL(ctx context.Context, short string) (string, error) {
	link, err := s.storage.Get(ctx, short)
	if err != nil {
		return "", err
	}

	if link.DeletedFlag {
		return "", ErrLinkIsDeleted
	}

	return link.OriginalURL.String(), nil
}

func (s ShortenerService) BatchCreateShortURL(ctx context.Context, urls []URLBatchInstance, userID int) ([]model.Link, error) {
	var links []model.Link
	for _, instance := range urls {
		short := s.generator.generateShort()
		links = append(links, model.NewLinkWithUUID(instance.Original, short, userID, instance.CorrelationID))
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

func (s ShortenerService) GetByUser(ctx context.Context, userID int) ([]model.Link, error) {
	return s.storage.GetByUser(ctx, userID)
}

func (s ShortenerService) DeleteByUser(ctx context.Context, userID int, shortURLs []string) error {
	s.deleteJobsChannel <- deleteTask{
		userID:    userID,
		shortURLs: shortURLs,
	}

	return nil
}

func (s ShortenerService) runDelete(ctx context.Context) error {
	if len(s.deleteJobsChannel) == 0 {
		return nil
	}

	jobsByUser := make(map[int][]string)
	for range len(s.deleteJobsChannel) {
		task := <-s.deleteJobsChannel
		jobsByUser[task.userID] = append(jobsByUser[task.userID], task.shortURLs...)
	}

	for userID, shortURLs := range jobsByUser {
		err := s.storage.DeleteByUser(ctx, userID, shortURLs)
		if err != nil {
			logger.Log.Info("Delete task failed", zap.Int("userID", userID), zap.Strings("shortURLs", shortURLs))
		}
	}

	return nil
}

func (s ShortenerService) RunDeleteScheduler(ctx context.Context) {
	logger.Log.Info("Delete queue started")
	for range s.deleteTicker.C {
		go s.runDelete(ctx)
	}
}

func (s ShortenerService) Close() {
	logger.Log.Info("Delete queue started")
	close(s.deleteJobsChannel)
	s.deleteTicker.Stop()
}

func NewShortenerService(repository repository.Repository, deleteInterval time.Duration) ShortenerService {
	return ShortenerService{
		storage:           repository,
		generator:         randomGenerator{},
		deleteJobsChannel: make(chan deleteTask, 100),
		deleteTicker:      time.NewTicker(deleteInterval),
	}
}
