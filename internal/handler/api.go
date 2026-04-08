package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/model"
	"frrenzy/urlshortener/internal/util/logger"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type apiHandler struct {
	Services
}

func (s apiHandler) createShortJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		logger.Log.Info(errContentType.Error(), zap.Error(errContentType))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errContentType.Error()))
		return
	}

	dec := json.NewDecoder(r.Body)
	var req model.ShortenRequest
	if err := dec.Decode(&req); err != nil {
		logger.Log.Info(errCanNotDecode.Error(), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errCanNotDecode.Error()))
		return
	}

	original, err := url.ParseRequestURI(req.URL)
	if err != nil {
		logger.Log.Info(errCanNotParseURL.Error(), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errCanNotParseURL.Error()))
		return
	}

	short, err := s.UrlService.CreateShortURL(*original)
	if err != nil {
		logger.Log.Info(errCanNotCreate.Error(), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errCanNotCreate.Error()))
		return
	}

	res := model.ShortenResponse{
		Result: config.Config.BaseAddress + "/" + short,
	}
	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusCreated)

	if err := enc.Encode(res); err != nil {
		logger.Log.Debug(errCanNotEncode.Error(), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func newAPIRouter(services Services) chi.Router {
	router := chi.NewRouter()
	apiHandlers := apiHandler{Services: services}

	router.Route(`/`, func(r chi.Router) {
		r.Post(`/shorten`, apiHandlers.createShortJSON)
	})

	return router
}
