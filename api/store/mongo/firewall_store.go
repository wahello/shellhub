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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Store) FirewallRuleList(ctx context.Context, pagination paginator.Query) ([]models.FirewallRule, int, error) {
	query := []bson.M{
		{
			"$sort": bson.M{
				"priority": 1,
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
	count, err := aggregateCount(ctx, s.db.Collection("firewall_rules"), queryCount)
	if err != nil {
		return nil, 0, fromMongoError(err)
	}

	query = append(query, buildPaginationQuery(pagination)...)

	rules := make([]models.FirewallRule, 0)
	cursor, err := s.db.Collection("firewall_rules").Aggregate(ctx, query)
	if err != nil {
		return nil, 0, fromMongoError(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		rule := new(models.FirewallRule)
		err = cursor.Decode(&rule)
		if err != nil {
			return rules, count, fromMongoError(err)
		}

		rules = append(rules, *rule)
	}

	return rules, count, fromMongoError(err)
}

func (s *Store) FirewallRuleCreate(ctx context.Context, rule *models.FirewallRule) error {
	if err := rule.Validate(); err != nil {
		return fromMongoError(err)
	}

	result, err := s.db.Collection("firewall_rules").InsertOne(ctx, &rule)
	if err != nil {
		return fromMongoError(err)
	}

	if err := s.cache.Set(ctx, strings.Join([]string{"rule", result.InsertedID.(primitive.ObjectID).String()}, "/"), rule, time.Minute); err != nil {
		logrus.Error(err)
	}

	return nil
}

func (s *Store) FirewallRuleGet(ctx context.Context, id string) (*models.FirewallRule, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fromMongoError(err)
	}

	var rule *models.FirewallRule
	if err := s.cache.Get(ctx, strings.Join([]string{"rule", id}, "/"), &rule); err != nil {
		logrus.Error(err)
	}
	if rule != nil {
		return rule, nil
	}

	if err := s.db.Collection("firewall_rules").FindOne(ctx, bson.M{"_id": objID}).Decode(&rule); err != nil {
		return nil, fromMongoError(err)
	}

	return rule, nil
}

func (s *Store) FirewallRuleUpdate(ctx context.Context, id string, rule models.FirewallRuleUpdate) (*models.FirewallRule, error) {
	if err := rule.Validate(); err != nil {
		return nil, fromMongoError(err)
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fromMongoError(err)
	}

	if _, err := s.db.Collection("firewall_rules").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": rule}); err != nil {
		return nil, fromMongoError(err)
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"rule", id}, "/")); err != nil {
		logrus.Error(err)
	}

	r, err := s.FirewallRuleGet(ctx, id)

	return r, fromMongoError(err)
}

func (s *Store) FirewallRuleDelete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fromMongoError(err)
	}

	if _, err := s.db.Collection("firewall_rules").DeleteOne(ctx, bson.M{"_id": objID}); err != nil {
		return fromMongoError(err)
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"rule", id}, "/")); err != nil {
		logrus.Error(err)
	}

	return nil
}
