package db

import "errors"

var (
	errCouldNotCreateDriver   = errors.New("could not create driver instance")
	errCouldNotCreateMigrator = errors.New("could not create migrator instance")
	errCouldNotRunMigrations  = errors.New("could not create run migrations")
)
