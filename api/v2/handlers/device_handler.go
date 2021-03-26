package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	apierrors "github.com/shellhub-io/shellhub/api/v2/errors"
	"github.com/shellhub-io/shellhub/api/v2/pkg/models"
)

func (h *Handler) registerDeviceRoutes() {
	h.external.GET(
		"/devices/list",
		h.deviceList,
	)
	h.external.POST(
		"/devices/list",
		h.deviceList,
	)
}

func (h *Handler) deviceList(c echo.Context) error {
	var params models.ListParams

	if err := c.Bind(&params); err != nil {
		return err
	}

	if err := params.IsValid(); err != nil {
		return apierrors.Wrap(err, apierrors.ErrInvalidArgument, "")
	}

	params.Pagination.ApplyLimits()

	ctx := c.Request().Context()

	devices, count, err := h.app.DeviceList(ctx, &params)
	if err != nil {
		return err
	}

	c.Response().Header().Set("X-Total-Count", strconv.Itoa(count))

	return c.JSON(http.StatusOK, devices)
}
