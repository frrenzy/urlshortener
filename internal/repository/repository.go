// Package repository
package repository

import (
	"frrenzy/urlshortener/internal/model"
)

type Repository interface {
	Add(l model.Link) error
	Get(short string) (model.Link, error)
}
