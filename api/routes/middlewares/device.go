package middlewares

import (
	"context"
	//	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/api/pkg/namespace"
	//	"errors"
	//	"fmt"
)

func AuthorizeDeviceOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), "ctx", c.(*apicontext.Context)) //nolint:revive
		tenant := ""
		if v := c.Tenant(); v != nil {
			tenant = v.ID
		}

		id := ""
		if v := c.ID(); v != nil {
			id = v.ID
		}

		err := namespace.IsNamespaceOwner(c.Ctx(), c.Store(), tenant, id)
		return c.NoContent(err)

	}

	return next(c)
}
