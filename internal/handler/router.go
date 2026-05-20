// Package handler
package handler

import (
	"net/http"

	"frrenzy/urlshortener/internal/service"

	"github.com/go-chi/chi/v5"
)

type Services struct {
	URLService service.ShortenerService
}

func NewRouter(services Services, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middlewares...)

	router.Mount(`/`, newPlainRouter(services))
	router.Mount(`/api`, newAPIRouter(services))
	router.Mount(`/ping`, newPingRouter(services))

	return router
}
