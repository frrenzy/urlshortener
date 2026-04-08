// Package config
package config

import "flag"

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseAddress   string `env:"BASE_URL"`
}

var Config config

func InitConfig() {
	flag.Parse()

	if Config.ServerAddress == "" {
		Config.ServerAddress = programFlags.ServerAddress
	}
	if Config.BaseAddress == "" {
		Config.BaseAddress = programFlags.BaseAddress
	}
}
