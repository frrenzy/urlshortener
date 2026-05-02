package handler

import (
	"io"
	"net/http"
	"net/url"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/util/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type plainHandler struct {
	Services
}

func (s plainHandler) createShort(w http.ResponseWriter, r *http.Request) {
	urlBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Info(errCanNotDecode.Error(), zap.Error(errContentType))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errCanNotDecode.Error()))
		return
	}

	original, err := url.ParseRequestURI(string(urlBytes))
	if err != nil {
		logger.Log.Info(errCanNotParseURL.Error(), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errCanNotParseURL.Error()))
		return
	}

	short, err := s.URLService.CreateShortURL(r.Context(), *original)
	if err != nil {
		logger.Log.Info(errCanNotCreate.Error(), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errCanNotCreate.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(config.Config.BaseAddress + "/" + short))
}

func (s plainHandler) redirectToOriginal(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "id")
	if short == "" {
		logger.Log.Info(errNoShortURLProvided.Error(), zap.Error(errNoShortURLProvided))
		w.Write([]byte(errNoShortURLProvided.Error()))
		return
	}

	original, err := s.URLService.GetOriginalURL(r.Context(), short)
	if err != nil {
		logger.Log.Info(errNoShortURLFound.Error(), zap.Error(errNoShortURLFound))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errNoShortURLFound.Error()))
		return
	}

	w.Header().Add("Location", original)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func newPlainRouter(services Services) chi.Router {
	router := chi.NewRouter()
	plainHandlers := plainHandler{Services: services}

	router.Route(`/`, func(r chi.Router) {
		r.Get(`/{id}`, plainHandlers.redirectToOriginal)
		r.Post(`/`, plainHandlers.createShort)
	})

	return router
}
