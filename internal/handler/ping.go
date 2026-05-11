package handler

import (
	"context"
	"net/http"
	"time"

	"frrenzy/urlshortener/internal/util/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type pingHandler struct {
	Services
}

func (s pingHandler) pingDB(w http.ResponseWriter, r *http.Request) {
	cancelCtx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	err := s.URLService.PingStorage(cancelCtx)
	if err != nil {
		logger.Log.Info("no connection", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func newPingRouter(services Services) chi.Router {
	router := chi.NewRouter()
	pingHandlers := pingHandler{Services: services}

	router.Route(`/`, func(r chi.Router) {
		r.Get(`/`, pingHandlers.pingDB)
	})

	return router
}
