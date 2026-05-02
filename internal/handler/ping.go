package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type pingHandler struct {
	Services
}

func (s pingHandler) pingDB(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cancelCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		err := s.DB.PingContext(cancelCtx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func newPingRouter(ctx context.Context, services Services) chi.Router {
	router := chi.NewRouter()
	pingHandlers := pingHandler{Services: services}

	router.Route(`/`, func(r chi.Router) {
		r.Get(`/`, pingHandlers.pingDB(ctx))
	})

	return router
}
