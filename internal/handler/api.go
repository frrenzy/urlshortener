package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/model"
	"frrenzy/urlshortener/internal/service"
	"frrenzy/urlshortener/internal/util/logger"
	"frrenzy/urlshortener/internal/util/user"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type apiHandler struct {
	Services
}

func (s apiHandler) createShortJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		logger.Log.Debug(errContentType.Error(), zap.Error(errContentType))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errContentType.Error()))
		return
	}

	dec := json.NewDecoder(r.Body)
	var req model.ShortenRequest
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug(errCanNotDecode.Error(), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errCanNotDecode.Error()))
		return
	}

	original, err := url.ParseRequestURI(req.URL)
	if err != nil {
		logger.Log.Debug(errCanNotParseURL.Error(), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errCanNotParseURL.Error()))
		return
	}

	ctx := r.Context()
	userID, err := user.GetUserFromContext(ctx)
	if err != nil {
		logger.Log.Info("context error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var statusCode int
	short, err := s.URLService.CreateShortURL(r.Context(), *original, userID)
	if err != nil {
		if !errors.Is(err, service.ErrLinkAlreadyExists) {
			logger.Log.Info(errCanNotCreate.Error(), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errCanNotCreate.Error()))

			return
		}

		logger.Log.Info(err.Error(), zap.Error(err))
		statusCode = http.StatusConflict
	} else {
		statusCode = http.StatusCreated
	}

	path, err := url.JoinPath(config.Config.BaseAddress, "/", short)
	if err != nil {
		logger.Log.Info(errCanNotCreate.Error(), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errCanNotCreate.Error()))

		return
	}

	res := model.ShortenResponse{
		Result: path,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)

	if err := enc.Encode(res); err != nil {
		logger.Log.Debug(errCanNotEncode.Error(), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s apiHandler) createShortJSONBatch(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		logger.Log.Debug(errContentType.Error(), zap.Error(errContentType))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errContentType.Error()))
		return
	}

	dec := json.NewDecoder(r.Body)
	var req model.BatchShortenRequest
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug(errCanNotDecode.Error(), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errCanNotDecode.Error()))
		return
	}

	var urls []service.URLBatchInstance
	for _, instance := range req {
		original, err := url.ParseRequestURI(instance.OriginalURL)
		if err != nil {
			logger.Log.Debug(errCanNotParseURL.Error(), zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errCanNotParseURL.Error()))
			return
		}
		urls = append(urls, service.URLBatchInstance{
			Original:      *original,
			CorrelationID: instance.CorrelationID,
		})
	}

	ctx := r.Context()
	userID, err := user.GetUserFromContext(ctx)
	if err != nil {
		logger.Log.Info("context error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	links, err := s.URLService.BatchCreateShortURL(r.Context(), urls, userID)
	if err != nil {
		logger.Log.Info(errCanNotCreate.Error(), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errCanNotCreate.Error()))
		return
	}

	var res model.BatchShortenResponse
	for _, link := range links {
		res = append(res, model.BatchShortenResponseInstance{
			CorrelationID: link.UUID,
			ShortURL:      config.Config.BaseAddress + "/" + link.ShortURL,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusCreated)

	if err := enc.Encode(res); err != nil {
		logger.Log.Debug(errCanNotEncode.Error(), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s apiHandler) getAllUserURLs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	userID, err := user.GetUserFromContext(ctx)
	if err != nil {
		logger.Log.Info("context error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	links, err := s.URLService.GetByUser(ctx, userID)
	if err != nil {
		logger.Log.Info("db error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(links) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var res model.UserURLsResponce
	for _, link := range links {
		res = append(res, model.UserURLsResponceInstance{
			OriginalURL: link.OriginalURL,
			ShortURL:    config.Config.BaseAddress + "/" + link.ShortURL,
		})
	}

	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)

	if err := enc.Encode(res); err != nil {
		logger.Log.Debug(errCanNotEncode.Error(), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func newAPIRouter(services Services) chi.Router {
	router := chi.NewRouter()
	apiHandlers := apiHandler{Services: services}

	router.Route(`/`, func(r chi.Router) {
		r.Post(`/shorten`, apiHandlers.createShortJSON)
		r.Post(`/shorten/batch`, apiHandlers.createShortJSONBatch)
		r.Get(`/user/urls`, apiHandlers.getAllUserURLs)
	})

	return router
}
