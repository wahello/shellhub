package namespace

import (
	"context"
	"errors"
	"fmt"

	"github.com/shellhub-io/shellhub/api/store"
)

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrUserNotFound      = errors.New("user not found")
	ErrNamespaceNotFound = errors.New("namespace not found")
)

func contains(members []interface{}, user string) bool {
	for _, member := range members {
		if member.(string) == user {
			return true
		}
	}
	return false
}

func IsNamespaceOwner(ctx context.Context, s store.Store, tenantID, ownerID string) error {
	user, _, err := s.UserGetByID(ctx, ownerID, false)
	if err == store.ErrNoDocuments {
		fmt.Println("ERRO 1")
		return ErrUnauthorized
	}

	if err != nil {
		fmt.Println("ERRO 2")
		return err
	}

	ns, err := s.NamespaceGet(ctx, tenantID)
	if err == store.ErrNoDocuments {
		fmt.Println("ERRO 3")
		return ErrNamespaceNotFound
	}

	if err != nil {
		fmt.Println("ERRO 4")
		return err
	}

	if ns.Owner != user.ID {
		fmt.Println("ERRO 5")
		return ErrUnauthorized
	}
	fmt.Println("SEM ERRO")
	return nil
}

func IsNamespaceMember(ctx context.Context, s store.Store, tenantID, memberID string) error {
	ns, err := s.NamespaceGet(ctx, tenantID)
	if err == store.ErrNoDocuments {
		return ErrNamespaceNotFound
	}
	if !contains(ns.Members, memberID) {
		return ErrUnauthorized
	}

	return nil

}
