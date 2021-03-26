package app

import (
	"context"
)

type UserService interface {
	UserLogin(ctx context.Context) error
}

func (a *app) UserLogin() error {
	return nil
}
