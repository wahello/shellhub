package apierrors

import (
	"github.com/pkg/errors"
)

// Kind is the kind of error
type Kind string

// Error kinds
const (
	ErrUnknown          Kind = "unknown_error"
	ErrInvalidArgument  Kind = "invalid_argument"
	ErrNotFound         Kind = "not_found"
	ErrAlreadyExists    Kind = "already_exists"
	ErrPermissionDenied Kind = "permission_denied"
	ErrUnauthenticated  Kind = "unauthenticated"
)

// Error is an API errors
type Error struct {
	error
	kind Kind
}

// Kind returns the error kind
func (e *Error) Kind() Kind {
	return e.kind
}

func (e *Error) Unwrap() error {
	return errors.Unwrap(Unwrap(e.error))
}

func (e *Error) Error() string {
	return e.error.Error()
}

// New returns an error with the supplied kind
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

// Wrap returns an error annotating err with a kind and the supplied message
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

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// Same as Go's errors.Unwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
