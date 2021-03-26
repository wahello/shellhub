package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	apierrors "github.com/shellhub-io/shellhub/api/v2/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		debug    bool
		response *responseError
	}{
		{
			name:  "internal echo error",
			err:   echo.NewHTTPError(http.StatusBadRequest, "error").SetInternal(errors.New("error")),
			debug: true,
			response: &responseError{
				Code:    "invalid_argument",
				Message: "error",
			},
		},
		{
			name:  "api error",
			err:   apierrors.Wrap(errors.New(""), apierrors.ErrInvalidArgument, ""),
			debug: true,
			response: &responseError{
				Code: "invalid_argument",
			},
		},
		{
			name:  "unknown error",
			err:   errors.New("error"),
			debug: true,
			response: &responseError{
				Code:    "unknown_error",
				Message: "error",
			},
		},
		{
			name:  "unknown error without debug message",
			err:   errors.New("error"),
			debug: false,
			response: &responseError{
				Code: "unknown_error",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			e.Debug = tc.debug

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := &Handler{router: e}

			h.errorHandler(tc.err, c)

			res := responseError{}
			err := json.Unmarshal([]byte(rec.Body.String()), &res)

			if assert.NoError(t, err) {
				assert.Equal(t, tc.response, &res)
			}
		})
	}

}
