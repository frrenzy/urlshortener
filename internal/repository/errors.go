package repository

import "errors"

var (
	errNotFound       = errors.New("not found")
	errDb             = errors.New("database error")
	errDBNotConnected = errors.New("no db connection")
)
