package haljson

import "errors"

var (
	// ErrNoCurie is returned when a curied link was added without the associated curie
	ErrNoCurie = errors.New("Must add curie before adding a curied link")
)
