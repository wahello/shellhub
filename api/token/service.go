package token

import (
	"context"

	"github.com/shellhub-io/shellhub/api/store"
	"github.com/shellhub-io/shellhub/pkg/models"
)

type Service interface {
	ListToken(ctx context.Context, tenantID string) ([]models.Token, error)
	CreateToken(ctx context.Context, tenantID string) (*models.Token, error)
	GetToken(ctx context.Context, tenantID string, ID string) (*models.Token, error)
	DeleteToken(ctx context.Context, tenantID string, ID string) error
	UpdateToken(ctx context.Context, tenantID string, ID string, token *models.APITokenUpdate) error
}

type service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return &service{store}
}

func (s *service) ListToken(ctx context.Context, tenantID string) ([]models.Token, error) {
	return s.store.NamespaceListAPIToken(ctx, tenantID)
}

func (s *service) CreateToken(ctx context.Context, tenantID string) (*models.Token, error) {
	return s.store.NamespaceCreateAPIToken(ctx, tenantID)
}

func (s *service) GetToken(ctx context.Context, tenantID string, ID string) (*models.Token, error) {
	return s.store.NamespaceGetAPIToken(ctx, tenantID, ID)
}

func (s *service) DeleteToken(ctx context.Context, tenantID string, ID string) error {
	return s.store.NamespaceDeleteAPIToken(ctx, tenantID, ID)
}

func (s *service) UpdateToken(ctx context.Context, tenatID string, ID string, request *models.APITokenUpdate) error {
	return s.store.NamespaceUpdateAPIToken(ctx, tenatID, ID, request)
}
