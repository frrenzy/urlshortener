// Package handler
package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/service"
)

func createShort(w http.ResponseWriter, r *http.Request) {
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

	short, err := service.CreateShortURL(*original)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("can not create short URL"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://localhost:%d/%s", config.Port, short)
}

func redirectToOriginal(w http.ResponseWriter, r *http.Request) {
	short := r.PathValue("id")
	if short == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no short URL provided"))
		return
	}

	original, err := service.GetOriginalURL(short)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no short URL found"))
		return
	}

	w.Header().Add("Location", original.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(""))
}

var MainHandler *http.ServeMux

func init() {
	MainHandler = http.NewServeMux()
	MainHandler.HandleFunc(`POST /`, createShort)
	MainHandler.HandleFunc(`GET /{id}`, redirectToOriginal)
}
