package main

import (
	"net/http"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/handler"
	"frrenzy/urlshortener/internal/repository"
	"frrenzy/urlshortener/internal/service"
	"frrenzy/urlshortener/internal/util/gzip"
	"frrenzy/urlshortener/internal/util/logger"
)

func run() error {
	config.Initialize()
	logger.Initialize(logger.DebugLevel)

	storage := repository.NewFileStorage(config.Config.FileStoragePath)
	defer storage.Close()

	services := handler.Services{
		URLService: service.NewShortenerService(storage),
	}

	router := handler.NewRouter(services, logger.WithLogging, gzip.CreateGzipMiddleware([]string{"text/html", "application/json"}))

	logger.Log.Info("Server listening on " + config.Config.ServerAddress)
	return http.ListenAndServe(config.Config.ServerAddress, router)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
