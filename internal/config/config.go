// Package config
package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseAddress     string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

var Config config

func Initialize() {
	err := env.Parse(&Config)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	if Config.ServerAddress == "" {
		Config.ServerAddress = programFlags.ServerAddress
	}
	if Config.BaseAddress == "" {
		Config.BaseAddress = programFlags.BaseAddress
	}
	if Config.FileStoragePath == "" {
		Config.FileStoragePath = programFlags.FileStoragePath
	}
}
