package mongo

import (
	"context"
	"strings"
	"time"

	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Store) PublicKeyGet(ctx context.Context, fingerprint, tenant string) (*models.PublicKey, error) {
	var pubKey *models.PublicKey
	if tenant != "" {
		if err := s.cache.Get(ctx, strings.Join([]string{"key", fingerprint}, "/"), &pubKey); err != nil {
			logrus.Error(err)
		}
		if pubKey != nil && pubKey.TenantID == tenant {
			return pubKey, nil
		}
		if err := s.db.Collection("public_keys").FindOne(ctx, bson.M{"fingerprint": fingerprint, "tenant_id": tenant}).Decode(&pubKey); err != nil {
			return nil, fromMongoError(err)
		}
	} else {
		if err := s.db.Collection("public_keys").FindOne(ctx, bson.M{"fingerprint": fingerprint}).Decode(&pubKey); err != nil {
			return nil, fromMongoError(err)
		}
	}

	return pubKey, nil
}

func (s *Store) PublicKeyList(ctx context.Context, pagination paginator.Query) ([]models.PublicKey, int, error) {
	query := []bson.M{
		{
			"$sort": bson.M{
				"created_at": 1,
			},
		},
	}

	// Only match for the respective tenant if requested
	if tenant := apicontext.TenantFromContext(ctx); tenant != nil {
		query = append(query, bson.M{
			"$match": bson.M{
				"tenant_id": tenant.ID,
			},
		})
	}

	queryCount := append(query, bson.M{"$count": "count"})
	count, err := aggregateCount(ctx, s.db.Collection("public_keys"), queryCount)
	if err != nil {
		return nil, 0, err
	}

	query = append(query, buildPaginationQuery(pagination)...)

	list := make([]models.PublicKey, 0)
	cursor, err := s.db.Collection("public_keys").Aggregate(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		key := new(models.PublicKey)
		err = cursor.Decode(&key)
		if err != nil {
			return list, count, err
		}

		list = append(list, *key)
	}

	return list, count, err
}

func (s *Store) PublicKeyCreate(ctx context.Context, key *models.PublicKey) error {
	if err := key.Validate(); err != nil {
		return err
	}

	_, err := s.db.Collection("public_keys").InsertOne(ctx, key)
	if err != nil {
		return fromMongoError(err)
	}

	if err := s.cache.Set(ctx, strings.Join([]string{"key", key.Fingerprint}, "/"), key, time.Minute); err != nil {
		logrus.Error(err)
	}

	return nil
}

func (s *Store) PublicKeyUpdate(ctx context.Context, fingerprint, tenant string, key *models.PublicKeyUpdate) (*models.PublicKey, error) {
	if err := key.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.db.Collection("public_keys").UpdateOne(ctx, bson.M{"fingerprint": fingerprint}, bson.M{"$set": key}); err != nil {
		if err != nil {
			return nil, fromMongoError(err)
		}

		return nil, err
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"key", fingerprint}, "/")); err != nil {
		logrus.Error(err)
	}

	return s.PublicKeyGet(ctx, fingerprint, tenant)
}

func (s *Store) PublicKeyDelete(ctx context.Context, fingerprint, tenant string) error {
	if _, err := s.db.Collection("public_keys").DeleteOne(ctx, bson.M{"fingerprint": fingerprint, "tenant_id": tenant}); err != nil {
		return fromMongoError(err)
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"key", fingerprint}, "/")); err != nil {
		logrus.Error(err)
	}

	return nil
}
