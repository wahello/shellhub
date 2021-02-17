package migrations

import (
	"context"

	"github.com/sirupsen/logrus"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var migration27 = migrate.Migration{
	Version:     27,
	Description: "Create a new field on namespaces called api_tokens to store the tokens",
	Up: func(db *mongo.Database) error {
		_, err := db.Collection("namespaces").UpdateMany(context.TODO(), bson.M{}, bson.M{"$set": bson.M{"api_tokens": []models.Token{}}})
		return err
	},
	Down: func(db *mongo.Database) error {
		return nil
	},
}
