// Package repository
package repository

import (
	"context"

	"frrenzy/urlshortener/internal/model"
)

type Repository interface {
	Add(ctx context.Context, link model.Link) error
	Get(ctx context.Context, short string) (*model.Link, error)
}
