package db

import "errors"

var (
	errCouldNotEstablishConnection = errors.New("could not establish DB connection")
	errCouldNotCreateDriver        = errors.New("could not create driver instance")
	errCouldNotCreateMigrator      = errors.New("could not create migrator instance")
	errCouldNotRunMigrations       = errors.New("could not create run migrations")
)
