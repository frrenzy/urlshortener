package handler

import "errors"

var (
	errContentType        = errors.New("wrong Content-Type")
	errCanNotDecode       = errors.New("can not decode request")
	errCanNotParseURL     = errors.New("can not parse url")
	errUnauthorized       = errors.New("user unauthorized")
	errCanNotCreate       = errors.New("can not create short URL")
	errCanNotEncode       = errors.New("can not encode request")
	errNoShortURLProvided = errors.New("no short URL provided")
	errNoShortURLFound    = errors.New("short URL not found")
)
