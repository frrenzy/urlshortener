// Package repository
package repository

import (
	"errors"
	"net/url"
)

var storage map[string]url.URL

func Add(original url.URL, short string) bool {
	storage[short] = original
	return true
}

func Get(short string) (url.URL, error) {
	original, exists := storage[short]
	if !exists {
		return url.URL{}, errors.New("not found")
	}

	return original, nil
}

func init() {
	storage = make(map[string]url.URL)
}
