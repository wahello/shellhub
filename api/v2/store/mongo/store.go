package mongostore

import (
	"github.com/shellhub-io/shellhub/api/v2/store"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongostore struct {
	db *mongo.Database
	store.Store
}

var _ store.Store = &mongostore{}

func NewStore(db *mongo.Database) store.Store {
	return &mongostore{db: db}
}
