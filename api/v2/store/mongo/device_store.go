package mongostore

import (
	"context"

	"github.com/shellhub-io/shellhub/api/v2/pkg/apicontext"
	"github.com/shellhub-io/shellhub/api/v2/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *mongostore) DeviceList(ctx context.Context, params *models.ListParams) ([]*models.Device, int, error) {
	queryMatch, err := buildFilterQuery(params.Filters)
	if err != nil {
		return nil, 0, err
	}

	query := []bson.M{
		{

			"$lookup": bson.M{
				"from":         "connected_devices",
				"localField":   "uid",
				"foreignField": "uid",
				"as":           "online",
			},
		},
		{
			"$addFields": bson.M{
				"online": bson.M{"$anyElementTrue": []interface{}{"$online"}},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "namespaces",
				"localField":   "tenant_id",
				"foreignField": "tenant_id",
				"as":           "namespace",
			},
		},
		{
			"$addFields": bson.M{
				"namespace": "$namespace.name",
			},
		},
		{
			"$unwind": "$namespace",
		},
	}

	// Default sort
	query = append(query, bson.M{
		"$sort": bson.M{"last_seen": -1},
	})

	// Apply filters
	if len(queryMatch) > 0 {
		query = append(query, queryMatch...)
	}

	// Only match for the respective tenant
	if session := apicontext.GetSessionContext(ctx); session != nil {
		query = append(query, bson.M{
			"$match": bson.M{
				"tenant_id": session.TenantID,
			},
		})
	}

	queryCount := append(query, bson.M{"$count": "count"})
	count, err := aggregateCount(ctx, s.db.Collection("devices"), queryCount)
	if err != nil {
		return nil, 0, err
	}

	// Append pagination query
	query = append(query, buildPaginationQuery(params.Pagination)...)

	cursor, err := s.db.Collection("devices").Aggregate(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(ctx)

	devices := make([]*models.Device, count)

	for cursor.Next(ctx) {
		device := new(models.Device)
		if err = cursor.Decode(&device); err != nil {
			return nil, 0, err
		}

		devices = append(devices, device)
	}

	return devices, count, err
}
