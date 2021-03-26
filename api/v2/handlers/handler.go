package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/shellhub-io/shellhub/api/v2/app"
	"github.com/shellhub-io/shellhub/api/v2/pkg/apicontext"
)

type Handler struct {
	app app.App
	// Main router
	router *echo.Echo
	// External router
	external *echo.Group
	// Internal router
	internal *echo.Group
}

func NewHandler(a app.App, router *echo.Echo) *Handler {
	h := &Handler{app: a, router: router}

	h.router.HTTPErrorHandler = h.errorHandler
	h.router.Use(h.sessionMiddleware)

	h.internal = h.router.Group("/internal")
	h.external = h.router.Group("/api")

	h.registerNamespaceRoutes()
	h.registerDeviceRoutes()
	h.registerUserRoutes()

	return h
}

// sessionMiddleware set session context to current request
func (h *Handler) sessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("aeXXXXX")
		ctx := apicontext.SetSessionContext(c.Request().Context(), &apicontext.SessionContext{
			TenantID: "xxxx",
		})

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
