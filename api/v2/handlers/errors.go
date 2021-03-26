package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	apierrors "github.com/shellhub-io/shellhub/api/v2/errors"
)

// responseError is a HTTP error response
type responseError struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

// errorHandler is an handler used to inform when an error has occurred
func (h *Handler) errorHandler(err error, c echo.Context) {
	// convert generic error from echo into our custom error type
	if e, ok := err.(*echo.HTTPError); ok {
		kind := errorKindFromStatusCode(e.Code)
		err = apierrors.Wrap(e.Unwrap(), kind, "")
	}

	// handle custom errors
	if e, ok := err.(*apierrors.Error); ok {
		if !c.Response().Committed {
			re := responseError{
				Code: string(e.Kind()),
			}

			if h.router.Debug {
				re.Message = e.Unwrap().Error()
			}

			c.JSON(statusCodeFromErrorKind(e.Kind()), re)
		}
	} else {
		h.errorHandler(apierrors.Wrap(err, apierrors.ErrUnknown, ""), c)
	}
}

// errorKindFromStatusCode converts a HTTP status code into error kind
func errorKindFromStatusCode(status int) apierrors.Kind {
	switch status {
	case http.StatusInternalServerError:
		return apierrors.ErrUnknown
	case http.StatusBadRequest:
		return apierrors.ErrInvalidArgument
	case http.StatusNotFound:
		return apierrors.ErrNotFound
	case http.StatusConflict:
		return apierrors.ErrAlreadyExists
	case http.StatusForbidden:
		return apierrors.ErrPermissionDenied
	case http.StatusUnauthorized:
		return apierrors.ErrUnauthenticated
	}

	return apierrors.ErrUnknown

}

// statusCodeFromErrorKind converts an error kind into HTTP status code
func statusCodeFromErrorKind(kind apierrors.Kind) int {
	switch kind {
	case apierrors.ErrUnknown:
		return http.StatusInternalServerError
	case apierrors.ErrInvalidArgument:
		return http.StatusBadRequest
	case apierrors.ErrNotFound:
		return http.StatusNotFound
	case apierrors.ErrAlreadyExists:
		return http.StatusConflict
	case apierrors.ErrPermissionDenied:
		return http.StatusForbidden
	case apierrors.ErrUnauthenticated:
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}
