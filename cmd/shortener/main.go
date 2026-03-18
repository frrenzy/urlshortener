package main

import (
	"fmt"
	"net/http"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/handler"
)

func run() error {
	fmt.Println("Server listening on localhost:8080")
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handler.MainHandler)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
