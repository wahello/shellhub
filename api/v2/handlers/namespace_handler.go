package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) registerNamespaceRoutes() {
	h.external.GET(
		"/namespaces/list",
		h.namespaceList,
		supportedAuthorizations("user", "app"),
	)
}

func (h *Handler) namespaceList(c echo.Context) error {
	return h.app.NamespaceList(c.Request().Context(), nil, nil)
}

// supportedAuthorizations middleware ensures that current context is authorized using at least one of authTypes
func supportedAuthorizations(types ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, v := range types {
				if v == c.Request().Header.Get("X-Authentication-Type") {
					return next(c)
				}
			}

			return c.NoContent(http.StatusForbidden)
		}
	}
}
