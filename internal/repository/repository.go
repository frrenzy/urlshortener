package repository

import "net/url"

type Repository interface {
	Add(original url.URL, short string)
	Get(short string) (url.URL, error)
}
