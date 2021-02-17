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

func (s *Store) NamespaceList(ctx context.Context, pagination paginator.Query, filters []models.Filter, export bool) ([]models.Namespace, int, error) {
	query := []bson.M{}
	queryMatch, err := buildFilterQuery(filters)
	if err != nil {
		return nil, 0, fromMongoError(err)
	}

	if len(queryMatch) > 0 {
		query = append(query, queryMatch...)
	}

	if export {
		query = []bson.M{

			{
				"$lookup": bson.M{
					"from":         "devices",
					"localField":   "tenant_id",
					"foreignField": "tenant_id",
					"as":           "devices",
				},
			},
			{
				"$addFields": bson.M{
					"devices": bson.M{"$size": "$devices"},
				},
			},
			{
				"$lookup": bson.M{
					"from":         "sessions",
					"localField":   "devices.uid",
					"foreignField": "device_uid",
					"as":           "sessions",
				},
			},
			{
				"$addFields": bson.M{
					"sessions": bson.M{"$size": "$sessions"},
				},
			},
		}
	}

	if len(queryMatch) > 0 {
		query = append(query, queryMatch...)
	}

	// Only match for the respective tenant if requested
	if id := apicontext.IDFromContext(ctx); id != nil {
		user, _, err := s.UserGetByID(ctx, id.ID, false)
		if err != nil {
			return nil, 0, err
		}

		query = append(query, bson.M{
			"$match": bson.M{
				"members": user.ID,
			},
		})
	}

	queryCount := append(query, bson.M{"$count": "count"})
	count, err := aggregateCount(ctx, s.db.Collection("namespaces"), queryCount)
	if err != nil {
		return nil, 0, err
	}

	if pagination.Page != 0 && pagination.PerPage != 0 {
		query = append(query, buildPaginationQuery(pagination)...)
	}

	namespaces := make([]models.Namespace, 0)
	cursor, err := s.db.Collection("namespaces").Aggregate(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		namespace := new(models.Namespace)
		err = cursor.Decode(&namespace)
		if err != nil {
			return namespaces, count, err
		}

		countDevice, err := s.db.Collection("devices").CountDocuments(ctx, bson.M{"tenant_id": namespace.TenantID, "status": "accepted"})
		if err != nil {
			return namespaces, 0, err
		}

		namespace.DevicesCount = int(countDevice)

		namespaces = append(namespaces, *namespace)
	}

	return namespaces, count, err
}

func (s *Store) NamespaceGet(ctx context.Context, tenantID string) (*models.Namespace, error) {
	var ns *models.Namespace

	if err := s.cache.Get(ctx, strings.Join([]string{"namespace", tenantID}, "/"), &ns); err != nil {
		logrus.Error(err)
	}

	if ns != nil {
		goto count
	}

	if err := s.db.Collection("namespaces").FindOne(ctx, bson.M{"tenant_id": tenantID}).Decode(&ns); err != nil {
		return ns, fromMongoError(err)
	}

	if err := s.cache.Set(ctx, strings.Join([]string{"namespace", tenantID}, "/"), ns, time.Minute); err != nil {
		logrus.Error(err)
	}

count:
	countDevice, err := s.db.Collection("devices").CountDocuments(ctx, bson.M{"tenant_id": tenantID, "status": "accepted"})
	if err != nil {
		return nil, fromMongoError(err)
	}

	ns.DevicesCount = int(countDevice)

	return ns, nil
}

func (s *Store) NamespaceGetByName(ctx context.Context, name string) (*models.Namespace, error) {
	var ns *models.Namespace

	if err := s.cache.Get(ctx, strings.Join([]string{"namespace", name}, "/"), &ns); err != nil {
		logrus.Error(err)
	}

	if ns != nil {
		return ns, nil
	}

	if err := s.db.Collection("namespaces").FindOne(ctx, bson.M{"name": name}).Decode(&ns); err != nil {
		return nil, fromMongoError(err)
	}

	return ns, nil
}

func (s *Store) NamespaceCreate(ctx context.Context, namespace *models.Namespace) (*models.Namespace, error) {
	_, err := s.db.Collection("namespaces").InsertOne(ctx, namespace)

	return namespace, err
}

func (s *Store) NamespaceDelete(ctx context.Context, tenantID string) error {
	if _, err := s.db.Collection("namespaces").DeleteOne(ctx, bson.M{"tenant_id": tenantID}); err != nil {
		return fromMongoError(err)
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"namespace", tenantID}, "/")); err != nil {
		logrus.Error(err)
	}

	collections := []string{"devices", "sessions", "connected_devices", "firewall_rules", "public_keys", "recorded_sessions"}
	for _, collection := range collections {
		if _, err := s.db.Collection(collection).DeleteMany(ctx, bson.M{"tenant_id": tenantID}); err != nil {
			return fromMongoError(err)
		}
	}

	return nil
}

func (s *Store) NamespaceRename(ctx context.Context, tenantID, name string) (*models.Namespace, error) {
	if _, err := s.db.Collection("namespaces").UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$set": bson.M{"name": name}}); err != nil {
		return nil, fromMongoError(err)
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"namespace", tenantID}, "/")); err != nil {
		logrus.Error(err)
	}

	return s.NamespaceGet(ctx, tenantID)
}

func (s *Store) NamespaceUpdate(ctx context.Context, tenantID string, namespace *models.Namespace) error {
	if _, err := s.db.Collection("namespaces").UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$set": bson.M{"name": namespace.Name, "max_devices": namespace.MaxDevices, "settings.session_record": namespace.Settings.SessionRecord}}); err != nil {
		return fromMongoError(err)
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"namespace", tenantID}, "/")); err != nil {
		logrus.Error(err)
	}

	return nil
}

func (s *Store) NamespaceAddMember(ctx context.Context, tenantID, id string) (*models.Namespace, error) {
	result, err := s.db.Collection("namespaces").UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$addToSet": bson.M{"members": id}})
	if err != nil {
		return nil, fromMongoError(err)
	}

	if result.ModifiedCount == 0 {
		return nil, ErrDuplicateID
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"namespace", tenantID}, "/")); err != nil {
		logrus.Error(err)
	}

	return s.NamespaceGet(ctx, tenantID)
}

func (s *Store) NamespaceRemoveMember(ctx context.Context, tenantID, id string) (*models.Namespace, error) {
	result, err := s.db.Collection("namespaces").UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$pull": bson.M{"members": id}})
	if err != nil {
		return nil, fromMongoError(err)
	}
	if result.ModifiedCount == 0 {
		return nil, ErrUserNotFound
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"namespace", tenantID}, "/")); err != nil {
		logrus.Error(err)
	}

	return s.NamespaceGet(ctx, tenantID)
}

func (s *Store) NamespaceGetFirst(ctx context.Context, id string) (*models.Namespace, error) {
	ns := new(models.Namespace)
	if err := s.db.Collection("namespaces").FindOne(ctx, bson.M{"members": id}).Decode(&ns); err != nil {
		return nil, fromMongoError(err)
	}

	return ns, nil
}

func (s *Store) NamespaceSetSessionRecord(ctx context.Context, sessionRecord bool, tenantID string) error {
	if _, err := s.db.Collection("namespaces").UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$set": bson.M{"settings.session_record": sessionRecord}}); err != nil {
		return fromMongoError(err)
	}

	return nil
}

func (s *Store) NamespaceGetSessionRecord(ctx context.Context, tenantID string) (bool, error) {
	var settings struct {
		Settings *models.NamespaceSettings `json:"settings" bson:"settings"`
	}

	if err := s.db.Collection("namespaces").FindOne(ctx, bson.M{"tenant_id": tenantID}).Decode(&settings); err != nil {
		return false, fromMongoError(err)
	}

	return settings.Settings.SessionRecord, nil
}

func (s *Store) NamespaceListAPIToken(ctx context.Context, tenantID string) ([]models.Token, error) {
	ns := new(models.Namespace)

	if err := s.db.Collection("namespaces").FindOne(ctx, bson.M{"tenant_id": tenantID}).Decode(&ns); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, store.ErrNamespaceNoDocuments
		}
		return nil, err
	}

	return ns.APITokens, nil
}

func (s *Store) NamespaceCreateAPIToken(ctx context.Context, tenantID string) (*models.Token, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(tenantID), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hasher := md5.New()
	if _, err := hasher.Write(hash); err != nil {
		return nil, err
	}

	token := &models.Token{
		ID:       hex.EncodeToString(hasher.Sum(nil)),
		TenantID: tenantID,
		ReadOnly: true,
	}

	_, err = s.db.Collection("namespaces").UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$push": bson.M{"api_tokens": token}})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Store) NamespaceGetAPIToken(ctx context.Context, tenantID string, ID string) (*models.Token, error) {
	tokens, err := s.ListAPIToken(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		if token.ID == ID {
			return &token, nil
		}
	}

	return nil, store.ErrRecordNotFound
}

func (s *Store) NamespaceDeleteAPIToken(ctx context.Context, tenantID string, ID string) error {
	tokens, err := s.ListAPIToken(ctx, tenantID)
	if err != nil {
		return err
	}

	for id, token := range tokens {
		if token.ID == ID {
			tokens = append(tokens[:id], tokens[id+1:]...)
		}
	}

	_, err = s.db.Collection("namespaces").UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$set": bson.M{"api_tokens": tokens}})
	if err != nil {
		return err
	}

	if err == mongo.ErrNoDocuments {
		return store.ErrRecordNotFound
	}

	return nil
}

func (s *Store) NamespaceUpdateAPIToken(ctx context.Context, tenantID string, ID string, request *models.APITokenUpdate) error {
	tokens, err := s.ListAPIToken(ctx, tenantID)
	if err != nil {
		return nil
	}

	for id, token := range tokens {
		if token.ID == ID {
			tokens[id].ReadOnly = request.TokenFields.ReadOnly
		}
	}

	_, err = s.db.Collection("namespaces").UpdateOne(ctx, bson.M{"tenant_id": tenantID}, bson.M{"$set": bson.M{"api_tokens": tokens}})
	if err != nil {
		return err
	}

	if err == mongo.ErrNoDocuments {
		return store.ErrRecordNotFound
	}

	return nil
}
