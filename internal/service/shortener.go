// Package service
package service

import (
	"errors"
	"net/url"

	"frrenzy/urlshortener/internal/repository"
)

func CreateShortURL(original url.URL) (string, error) {
	short := generateShort()

	if ok := repository.Add(original, short); !ok {
		return "", errors.New("error creating url")
	}

	return short, nil
}

func GetOriginalURL(short string) (url.URL, error) {
	original, err := repository.Get(short)
	if err != nil {
		return url.URL{}, errors.New("not found")
	}

	return original, nil
}
