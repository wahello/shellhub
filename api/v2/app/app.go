package app

import (
	"context"

	"github.com/shellhub-io/shellhub/api/v2/pkg/apicontext"
	"github.com/shellhub-io/shellhub/api/v2/store"
)

type App interface {
	SessionContext(ctx context.Context) *apicontext.SessionContext

	NamespaceService
	DeviceService
}

type app struct {
	store store.Store
}

var _ App = &app{}

func NewApp(store store.Store) App {
	return &app{store: store}
}

// SessionContext returns the current session context
func (a *app) SessionContext(ctx context.Context) *apicontext.SessionContext {
	return apicontext.GetSessionContext(ctx)
}
