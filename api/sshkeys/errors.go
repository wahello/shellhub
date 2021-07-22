package sshkeys

import (
	"errors"
)

var (
	ErrInvalidFormat        = errors.New("invalid format")
	ErrDuplicateFingerprint = errors.New("this fingerprint already exits")
	ErrUnauthorized         = errors.New("unauthorized")
)
