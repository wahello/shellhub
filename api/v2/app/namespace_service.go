package app

import (
	"context"

	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/models"
)

type NamespaceService interface {
	NamespaceList(ctx context.Context, pagination *paginator.Query, filters []*models.Filter) error
	NamespaceGet(ctx context.Context) error
}

func (a *app) NamespaceList(ctx context.Context, pagination *paginator.Query, filters []*models.Filter) error {
	err := a.store.NamespaceList(ctx)
	return err
}

func (a *app) NamespaceGet(ctxx context.Context) error {
	return nil
}
