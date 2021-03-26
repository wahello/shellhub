package handlers

import "github.com/labstack/echo/v4"

func (h *Handler) registerUserRoutes() {
	h.external.POST(
		"/user/login",
		h.userLogin,
	)
}

func (h *Handler) userLogin(c echo.Context) error {
	return nil
}
