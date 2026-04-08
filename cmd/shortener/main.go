package main

import (
	"fmt"
	"net/http"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/handler"
)

func run() error {
	config.InitConfig()

	fmt.Println("Server listening on ", config.Config.ServerAddress)
	return http.ListenAndServe(config.Config.ServerAddress, handler.NewRouter())
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
