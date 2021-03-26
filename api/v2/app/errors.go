package app

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// ErrorKind is the kind of error.
type ErrorKind string

// Error kinds
const (
	ErrUnknown          ErrorKind = "unknown error"
	ErrInvalidArgument  ErrorKind = "invalid argument"
	ErrNotFound         ErrorKind = "entity not found"
	ErrAlreadyExists    ErrorKind = "already exists"
	ErrPermissionDenied ErrorKind = "permission denied"
	ErrUnauthenticated  ErrorKind = "unauthenticated"
)

// Error is an internal errors with stacktrace. It can be converted to a HTTP response
type Error struct {
	error
	kind ErrorKind
}

// Format formats the error.
func (e *Error) Format(s fmt.State, verb rune) {
	if formatter, ok := e.error.(fmt.Formatter); ok {
		formatter.Format(s, verb)
		return
	}
	io.WriteString(s, e.Error())
}

// ErrorKind returns the error kind.
func (e *Error) ErrorKind() ErrorKind {
	return e.kind
}

func (e *Error) Unwrap() error {
	return errors.Unwrap(Unwrap(e.error))
}

func (e *Error) Error() string {
	return e.error.Error()
}

// New returns an error with the supplied kind and message. If message is empty, a default message
// for the error kind will be used.
func New(kind ErrorKind) error {
	return &Error{
		error: errors.New(""),
		kind:  kind,
	}
}

// New returns an error with the supplied kind and message. If message is empty, a default message
// for the error kind will be used.
func NewN(kind ErrorKind, msg string) error {
	if msg == "" {
		msg = string(kind)
	}
	return &Error{
		error: errors.New(msg),
		kind:  kind,
	}
}

// Errorf formats according to a format specifier and return an unknown error with the string.
func Errorf(kind ErrorKind, format string, args ...interface{}) error {
	return NewN(kind, fmt.Sprintf(format, args...))
}

// Wrap returns an error annotating err with a kind and a stacktrace at the point Wrap is called,
// and the supplied kind and message. If err is nil, Wrap returns nil.
func Wrap(err error, kind ErrorKind, msg string) error {
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

// Wrapf returns an error annotating err with a stack trace at the point Wrapf is called, and the
// kind and format specifier. If err is nil, Wrapf returns nil.
func Wrapf(err error, kind ErrorKind, format string, args ...interface{}) error {
	return Wrap(err, kind, fmt.Sprintf(format, args...))
}

// IsErrorKind checks whether any error in err's chain matches the error kind.
func IsErrorKind(err error, kind ErrorKind) bool {
	ie := &Error{}
	if As(err, &ie) {
		return ie.kind == kind
	}
	return false
}

// As finds the first error in err's chain that matches target, and if so, sets target to that
// error value and return true.
//
// Same as Go's errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// Same as Go's errors.Unwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
