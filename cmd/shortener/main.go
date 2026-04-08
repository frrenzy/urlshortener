package main

import (
	"fmt"
	"net/http"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/handler"
	"frrenzy/urlshortener/internal/logger"
)

func run() error {
	config.InitConfig()
	logger.Initialize()

	fmt.Println("Server listening on ", config.Config.ServerAddress)
	return http.ListenAndServe(config.Config.ServerAddress, logger.WithLogging(handler.NewRouter()))
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
