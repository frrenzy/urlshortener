package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func (s handler) createShort(w http.ResponseWriter, r *http.Request) {
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

	short, err := s.urlService.CreateShortURL(*original)
	if err != nil {
		logger.Log.Info(errCanNotCreate.Error(), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errCanNotCreate.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", config.Config.BaseAddress, short)
}

func (s handler) redirectToOriginal(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "id")
	if short == "" {
		logger.Log.Info(errNoShortURLProvided.Error(), zap.Error(errNoShortURLProvided))
		w.Write([]byte(errNoShortURLProvided.Error()))
		return
	}

	original, err := s.urlService.GetOriginalURL(short)
	if err != nil {
		logger.Log.Info(errNoShortURLFound.Error(), zap.Error(errNoShortURLFound))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errNoShortURLFound.Error()))
		return
	}

	w.Header().Add("Location", original.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}
