package model

import "errors"

var (
	errDBScan    = errors.New("cannot scan db value")
	errDBConvert = errors.New("cannot convert db value")
)
