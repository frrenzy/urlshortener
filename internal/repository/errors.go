package repository

import "errors"

var (
	errDBNotFound     = errors.New("not found")
	ErrDBExisting     = errors.New("already exists in db")
	errDBNotConnected = errors.New("no db connection")
)
