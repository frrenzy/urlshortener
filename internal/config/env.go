// Package config
package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

var envVars config

func init() {
	err := env.Parse(&envVars)
	if err != nil {
		log.Fatal(err)
	}
}
