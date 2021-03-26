package app

import (
	"context"

	"github.com/shellhub-io/shellhub/api/v2/pkg/models"
)

type DeviceService interface {
	DeviceList(ctx context.Context, params *models.ListParams) ([]*models.Device, int, error)
}

func (a *app) DeviceList(ctx context.Context, params *models.ListParams) ([]*models.Device, int, error) {
	return a.store.DeviceList(ctx, params)
}
