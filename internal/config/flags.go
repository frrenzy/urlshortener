package config

import (
	"flag"
)

var programFlags config

const (
	defaultAddress = "localhost:8080"
)

func init() {
	flag.StringVar(&programFlags.ServerAddress, "a", defaultAddress, "server address")
	flag.StringVar(&programFlags.BaseAddress, "b", "http://"+defaultAddress, "base short URL address")
}
