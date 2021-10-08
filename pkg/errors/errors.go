package errors

import (
	"errors"
)

var (
	ErrReport         = errors.New("report error")
	ErrLocked = errors.New("pending request")
)
