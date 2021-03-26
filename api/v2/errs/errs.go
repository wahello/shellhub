package errs

import (
	"fmt"

	"github.com/pkg/errors"
)

// Kind is the kind of error.
type Kind string

// Error kinds
const (
	ErrUnknown          Kind = "unknown error"
	ErrInvalidArgument  Kind = "invalid argument"
	ErrNotFound         Kind = "entity not found"
	ErrAlreadyExists    Kind = "already exists"
	ErrPermissionDenied Kind = "permission denied"
	ErrUnauthenticated  Kind = "unauthenticated"
)

// Error is an internal error
type Error struct {
	error
	kind Kind
}

// Kind returns the error kind.
func (e *Error) Kind() Kind {
	return e.kind
}

func (e *Error) Unwrap() error {
	return errors.Unwrap(errors.Unwrap(e.error))
}

func (e *Error) Error() string {
	return e.error.Error()
}

// New returns an error with the supplied kind and empty message
func New(kind Kind) error {
	return &Error{
		error: errors.New(""),
		kind:  kind,
	}
}

// New returns an error with the supplied kind and message
func NewWithMessage(kind Kind, msg string) error {
	if msg == "" {
		msg = string(kind)
	}
	return &Error{
		error: errors.New(msg),
		kind:  kind,
	}
}

// Errorf formats according to a format specifier and return an unknown error with the string.
func Errorf(kind Kind, format string, args ...interface{}) error {
	return NewN(kind, fmt.Sprintf(format, args...))
}

// Wrap returns an error annotating err with a kind and supplied message
func Wrap(err error, kind Kind, msg string) error {
	if err == nil {
		return nil
	}
	if msg == "" {
		msg = string(kind)
	}
	return &Error{
		error: errors.Wrap(err, msg),
		kind:  kind,
	}
}
