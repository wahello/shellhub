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
	return s.store.TokenListAPIToken(ctx, tenantID)
}

func (s *service) CreateToken(ctx context.Context, tenantID string) (*models.Token, error) {
	return s.store.TokenCreateAPIToken(ctx, tenantID)
}

func (s *service) GetToken(ctx context.Context, tenantID string, id string) (*models.Token, error) {
	return s.store.TokenGetAPIToken(ctx, tenantID, id)
}

func (s *service) DeleteToken(ctx context.Context, tenantID string, id string) error {
	return s.store.TokenDeleteAPIToken(ctx, tenantID, id)
}

func (s *service) UpdateToken(ctx context.Context, tenatID string, id string, request *models.APITokenUpdate) error {
	return s.store.TokenUpdateAPIToken(ctx, tenatID, id, request)
}
