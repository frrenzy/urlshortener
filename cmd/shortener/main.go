package main

import (
	"context"
	"fmt"
	"net/http"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/config/db"
	"frrenzy/urlshortener/internal/handler"
	"frrenzy/urlshortener/internal/repository"
	"frrenzy/urlshortener/internal/service"
	"frrenzy/urlshortener/internal/util/gzip"
	"frrenzy/urlshortener/internal/util/logger"
)

func run() error {
	config.Initialize()
	logger.Initialize(logger.DebugLevel)

	fmt.Println(config.Config)

	storage := repository.NewFileStorage(config.Config.FileStoragePath)
	defer storage.Close()

	db, err := db.Connect(config.Config.DatabaseDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	services := handler.Services{
		URLService: service.NewShortenerService(storage),
		DB:         db,
	}

	ctx := context.Background()

	router := handler.NewRouter(ctx, services, logger.WithLogging, gzip.CreateGzipMiddleware([]string{"text/html", "application/json"}))

	logger.Log.Info("Server listening on " + config.Config.ServerAddress)
	return http.ListenAndServe(config.Config.ServerAddress, router)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
