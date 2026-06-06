// Package repository
package repository

import (
	"context"
	"net/url"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/config/db"
	"frrenzy/urlshortener/internal/model"
	"frrenzy/urlshortener/internal/util/logger"

	"go.uber.org/zap"
)

type Repository interface {
	Add(ctx context.Context, link model.Link) error
	BatchAdd(ctx context.Context, links []model.Link) ([]model.Link, error)
	Get(ctx context.Context, short string) (*model.Link, error)
	GetByOriginal(ctx context.Context, original url.URL) (*model.Link, error)
	GetByUser(ctx context.Context, userID int) ([]model.Link, error)
	DeleteByUser(ctx context.Context, userID int, shortURLs []string) error
	Close()
	PingStorage(ctx context.Context) error
}

func NewStorage() Repository {
	if config.Config.DatabaseDSN != "" {
		db, err := db.Connect(config.Config.DatabaseDSN)
		if err == nil {
			logger.Log.Debug("using DB storage")
			return NewDBStorage(db)
		} else {
			logger.Log.Debug("error db", zap.Error(err))
		}
	}
	if config.Config.FileStoragePath != "" {
		logger.Log.Debug("using file storage")
		return NewFileStorage(config.Config.FileStoragePath)
	}

	logger.Log.Debug("using map storage")
	return NewMapStorage()
}
