// Package handler
package handler

import (
	"frrenzy/urlshortener/internal/repository"
	"frrenzy/urlshortener/internal/service"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	urlService service.ShortenerService
}

func NewRouter() *chi.Mux {
	storage := repository.NewMapStorage()
	handlers := handler{
		urlService: service.NewShortenerService(storage),
	}

	router := chi.NewRouter()

	router.Get(`/{id}`, handlers.redirectToOriginal)
	router.Post(`/`, handlers.createShort)
	router.Post(`/api/shorten`, handlers.createShortJSON)

	return router
}
