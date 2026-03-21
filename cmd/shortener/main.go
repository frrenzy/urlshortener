package main

import (
	"flag"
	"fmt"
	"net/http"

	"frrenzy/urlshortener/internal/config"
	"frrenzy/urlshortener/internal/handler"
)

func run() error {
	flag.Parse()

	fmt.Println("Server listening on ", config.ProgramFlags.ServerAddress)
	return http.ListenAndServe(config.ProgramFlags.ServerAddress, handler.NewRouter())
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
