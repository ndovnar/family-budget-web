package store

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrDuplicateKey = errors.New("dublicate key")
	ErrUnknown      = errors.New("unknown")
)
