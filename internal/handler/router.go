// Package handler
package handler

import (
	"context"
	"database/sql"
	"net/http"

	"frrenzy/urlshortener/internal/service"

	"github.com/go-chi/chi/v5"
)

type Services struct {
	URLService service.ShortenerService
	DB         *sql.DB
}

func NewRouter(ctx context.Context, services Services, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middlewares...)

	router.Mount(`/`, newPlainRouter(ctx, services))
	router.Mount(`/api`, newAPIRouter(ctx, services))
	router.Mount(`/ping`, newPingRouter(ctx, services))

	return router
}
