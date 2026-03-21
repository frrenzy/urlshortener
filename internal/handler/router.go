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
)

type handler struct {
	urlService service.ShortenerService
}

func (s handler) createShort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong method"))
		return
	}

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
	fmt.Fprintf(w, "http://localhost:%d/%s", config.Port, short)
}

func (s handler) redirectToOriginal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong method"))
		return
	}

	short := r.PathValue("id")
	fmt.Println("short id = ", short)
	fmt.Println("request path = ", r.URL.Path)
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

func NewRouter() *http.ServeMux {
	storage := repository.NewMapStorage()
	mainRouter := handler{
		urlService: service.NewShortenerService(storage),
	}

	handler := http.NewServeMux()
	handler.HandleFunc(`/{id}`, mainRouter.redirectToOriginal)
	handler.HandleFunc(`/`, mainRouter.createShort)

	return handler
}
