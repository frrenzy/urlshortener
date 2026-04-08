// Package config
package config

import (
	"flag"
)

type flags struct {
	ServerAddress string
	BaseAddress   string
}

var ProgramFlags = flags{}

const (
	defaultAddress = "localhost:8080"
)

func init() {
	flag.StringVar(&ProgramFlags.ServerAddress, "a", defaultAddress, "server address")
	flag.StringVar(&ProgramFlags.BaseAddress, "b", "http://"+defaultAddress, "base short URL address")
}
