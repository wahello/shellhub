package store

import (
	"context"

	"github.com/shellhub-io/shellhub/api/v2/pkg/models"
)

type Store interface {
	DeviceStore
	NamespaceStore
}

type DeviceStore interface {
	DeviceList(ctx context.Context, params *models.ListParams) ([]*models.Device, int, error)
}

type NamespaceStore interface {
	NamespaceList(ctx context.Context) error
}
