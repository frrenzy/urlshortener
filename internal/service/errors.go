package service

import "errors"

var (
	ErrLinkAlreadyExists = errors.New("link already exists")
	ErrLinkIsDeleted     = errors.New("link is deleted")
)
