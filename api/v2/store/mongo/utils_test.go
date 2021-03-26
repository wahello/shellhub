package mongostore

import (
	"context"
	"testing"

	"github.com/shellhub-io/shellhub/api/pkg/dbtest"
	"github.com/shellhub-io/shellhub/api/v2/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestBuildFilterQueryInvalid(t *testing.T) {
	filter := models.FilterList{
		{
			Type: "invalid",
		},
	}

	query, err := buildFilterQuery(filter)

	assert.Nil(t, query)
	assert.EqualError(t, err, ErrInvalidFilter.Error())
}

func TestBuildFilterQueryWithComparasionOperators(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		value    string
		query    bson.M
	}{
		{
			name:     "contains operator",
			operator: "contains",
			value:    "value",
			query:    bson.M{"$regex": "value", "$options": "i"},
		},
		{
			name:     "eq operator",
			operator: "eq",
			value:    "value",
			query:    bson.M{"$eq": "value"},
		},
		{
			name:     "bool operator",
			operator: "bool",
			value:    "true",
			query:    bson.M{"$eq": true},
		},
		{
			name:     "gt operator",
			operator: "gt",
			value:    "8",
			query:    bson.M{"$gt": 8},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filter := models.FilterList{
				{
					Type: "property",
					Params: &models.FilterTypeProperty{
						Name:     "field",
						Operator: tc.operator,
						Value:    tc.value,
					},
				},
			}

			expectedQuery := []bson.M{
				bson.M{
					"$match": bson.M{
						"$or": []bson.M{
							bson.M{
								"field": tc.query,
							},
						},
					},
				},
			}

			query, err := buildFilterQuery(filter)

			assert.NoError(t, err)
			assert.NotNil(t, query)
			assert.Equal(t, expectedQuery, query)
		})
	}
}

func TestBuildPaginationQuery(t *testing.T) {
	tests := []struct {
		name       string
		pagination models.Pagination
		query      []bson.M
	}{
		{
			name:       "without limit",
			pagination: models.Pagination{PerPage: -1},
			query:      nil,
		},
		{
			name:       "with limit",
			pagination: models.Pagination{Page: 1, PerPage: 10},
			query: []bson.M{
				{"$skip": 0},
				{"$limit": 10},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			query := buildPaginationQuery(tc.pagination)

			assert.Equal(t, tc.query, query)
		})
	}
}

func TestAggregateCount(t *testing.T) {
	db := dbtest.DBServer{}
	defer db.Stop()

	ctx := context.TODO()
	coll := db.Client().Database("test").Collection("test")

	expectedCount := 5

	for i := 0; i < expectedCount; i++ {
		doc := map[string]interface{}{
			"key": "value",
		}

		_, err := coll.InsertOne(ctx, doc)
		assert.NoError(t, err)
	}

	count, err := aggregateCount(ctx, coll, []bson.M{{"$count": "count"}})

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}
