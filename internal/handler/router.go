// Package handler
package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/repository"
	"frrenzy/urlshortener/internal/service"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	urlService service.ShortenerService
}

func (s handler) createShort(w http.ResponseWriter, r *http.Request) {
	urlBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can not read request"))
		return
	}

	original, err := url.ParseRequestURI(string(urlBytes))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can not parse request"))
		return
	}

	short, err := s.urlService.CreateShortURL(*original)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can not create short URL"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", config.Config.BaseAddress, short)
}

func (s handler) redirectToOriginal(w http.ResponseWriter, r *http.Request) {
	short := chi.URLParam(r, "id")
	if short == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no short URL provided"))
		return
	}

	original, err := s.urlService.GetOriginalURL(short)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no short URL found"))
		return
	}

	w.Header().Add("Location", original.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func NewRouter() *chi.Mux {
	storage := repository.NewMapStorage()
	handlers := handler{
		urlService: service.NewShortenerService(storage),
	}

	router := chi.NewRouter()

	router.Get(`/{id}`, handlers.redirectToOriginal)
	router.Post(`/`, handlers.createShort)

	return router
}
