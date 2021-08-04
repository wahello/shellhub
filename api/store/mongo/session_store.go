package mongo

import (
	"context"
	"strings"
	"time"

	"github.com/shellhub-io/shellhub/api/apicontext"
	"github.com/shellhub-io/shellhub/pkg/api/paginator"
	"github.com/shellhub-io/shellhub/pkg/clock"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Store) SessionList(ctx context.Context, pagination paginator.Query) ([]models.Session, int, error) {
	query := []bson.M{
		{
			"$sort": bson.M{
				"started_at": -1,
			},
		},

		{
			"$lookup": bson.M{
				"from":         "active_sessions",
				"localField":   "uid",
				"foreignField": "uid",
				"as":           "active",
			},
		},
		{
			"$addFields": bson.M{
				"active": bson.M{"$anyElementTrue": []interface{}{"$active"}},
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
	count, err := aggregateCount(ctx, s.db.Collection("sessions"), queryCount)
	if err != nil {
		return nil, 0, fromMongoError(err)
	}

	query = append(query, buildPaginationQuery(pagination)...)

	sessions := make([]models.Session, 0)
	cursor, err := s.db.Collection("sessions").Aggregate(ctx, query)
	if err != nil {
		return sessions, count, fromMongoError(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		session := new(models.Session)
		err = cursor.Decode(&session)
		if err != nil {
			return sessions, count, err
		}

		device, err := s.DeviceGet(ctx, session.DeviceUID)
		if err != nil {
			return sessions, count, err
		}

		session.Device = device
		sessions = append(sessions, *session)
	}

	return sessions, count, err
}

func (s *Store) SessionGet(ctx context.Context, uid models.UID) (*models.Session, error) {
	var session *models.Session
	//how put cache here?
	query := []bson.M{
		{
			"$match": bson.M{"uid": uid},
		},
		{
			"$lookup": bson.M{
				"from":         "active_sessions",
				"localField":   "uid",
				"foreignField": "uid",
				"as":           "active",
			},
		},
		{
			"$addFields": bson.M{
				"active": bson.M{"$anyElementTrue": []interface{}{"$active"}},
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

	cursor, err := s.db.Collection("sessions").Aggregate(ctx, query)
	if err != nil {
		return nil, fromMongoError(err)
	}
	defer cursor.Close(ctx)
	cursor.Next(ctx)

	err = cursor.Decode(&session)
	if err != nil {
		return nil, fromMongoError(err)
	}

	device, err := s.DeviceGet(ctx, session.DeviceUID)
	if err != nil {
		return nil, fromMongoError(err)
	}

	session.Device = device

	return session, nil
}

func (s *Store) SessionSetAuthenticated(ctx context.Context, uid models.UID, authenticated bool) error {
	if _, err := s.db.Collection("sessions").UpdateOne(ctx, bson.M{"uid": uid}, bson.M{"$set": bson.M{"authenticated": authenticated}}); err != nil {
		return fromMongoError(err)

	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"session", string(uid)}, "/")); err != nil {
		logrus.Error(err)
	}

	return nil
}

func (s *Store) SessionSetRecorded(ctx context.Context, uid models.UID, recorded bool) error {
	if _, err := s.db.Collection("sessions").UpdateOne(ctx, bson.M{"uid": uid}, bson.M{"$set": bson.M{"recorded": recorded}}); err != nil {
		return fromMongoError(err)

	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"session", string(uid)}, "/")); err != nil {
		logrus.Error(err)
	}

	return nil

}

func (s *Store) SessionCreate(ctx context.Context, session models.Session) (*models.Session, error) {
	session.StartedAt = clock.Now()
	session.LastSeen = session.StartedAt
	session.Recorded = false

	device, err := s.DeviceGet(ctx, session.DeviceUID)
	if err != nil {
		return nil, fromMongoError(err)
	}

	session.TenantID = device.TenantID

	if _, err := s.db.Collection("sessions").InsertOne(ctx, &session); err != nil {
		return nil, fromMongoError(err)
	}

	if err := s.cache.Set(ctx, strings.Join([]string{"session", session.UID}, "/"), session, time.Minute); err != nil {
		logrus.Error(err)
	}

	as := &models.ActiveSession{
		UID:      models.UID(session.UID),
		LastSeen: session.StartedAt,
	}

	if _, err := s.db.Collection("active_sessions").InsertOne(ctx, &as); err != nil {
		return nil, fromMongoError(err)
	}

	return &session, nil
}

func (s *Store) SessionSetLastSeen(ctx context.Context, uid models.UID) error {
	var session *models.Session

	if err := s.cache.Get(ctx, strings.Join([]string{"session", string(uid)}, "/"), &session); err != nil {
		logrus.Error(err)
	}
	if session == nil {

		err := s.db.Collection("sessions").FindOne(ctx, bson.M{"uid": uid}).Decode(&session)
		if err != nil {
			return fromMongoError(err)
		}
	}

	if session != nil && session.Closed {
		return nil
	}

	session.LastSeen = clock.Now()

	opts := options.Update().SetUpsert(true)
	if _, err := s.db.Collection("sessions").UpdateOne(ctx, bson.M{"uid": session.UID}, bson.M{"$set": session}, opts); err != nil {
		return fromMongoError(err)
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"session", string(uid)}, "/")); err != nil {
		logrus.Error(err)
	}

	activeSession := &models.ActiveSession{
		UID:      uid,
		LastSeen: clock.Now(),
	}

	if _, err := s.db.Collection("active_sessions").InsertOne(ctx, &activeSession); err != nil {
		return fromMongoError(err)
	}

	return nil
}

func (s *Store) SessionDeleteActives(ctx context.Context, uid models.UID) error {
	var session *models.Session
	if err := s.db.Collection("sessions").FindOne(ctx, bson.M{"uid": uid}).Decode(&session); err != nil {
		return fromMongoError(err)
	}

	session.LastSeen = clock.Now()
	session.Closed = true

	opts := options.Update().SetUpsert(true)
	_, err := s.db.Collection("sessions").UpdateOne(ctx, bson.M{"uid": session.UID}, bson.M{"$set": session}, opts)
	if err != nil {
		return fromMongoError(err)
	}

	if err := s.cache.Delete(ctx, strings.Join([]string{"session", string(uid)}, "/")); err != nil {
		logrus.Error(err)
	}

	_, err = s.db.Collection("active_sessions").DeleteMany(ctx, bson.M{"uid": session.UID})

	return fromMongoError(err)
}

func (s *Store) SessionCreateRecordFrame(ctx context.Context, uid models.UID, recordSession *models.RecordedSession) error {
	if _, err := s.db.Collection("recorded_sessions").InsertOne(ctx, &recordSession); err != nil {
		return fromMongoError(err)
	}

	if _, err := s.db.Collection("sessions").UpdateOne(ctx, bson.M{"uid": uid}, bson.M{"$set": bson.M{"recorded": true}}); err != nil {
		return fromMongoError(err)
	}

	return nil
}

func (s *Store) SessionUpdateDeviceUID(ctx context.Context, oldUID models.UID, newUID models.UID) error {
	_, err := s.db.Collection("sessions").UpdateMany(ctx, bson.M{"device_uid": oldUID}, bson.M{"$set": bson.M{"device_uid": newUID}})

	//how put cache here?
	return fromMongoError(err)
}

func (s *Store) SessionDeleteRecordFrame(ctx context.Context, uid models.UID) error {
	_, err := s.db.Collection("recorded_sessions").DeleteMany(ctx, bson.M{"uid": uid})

	return fromMongoError(err)
}

func (s *Store) SessionGetRecordFrame(ctx context.Context, uid models.UID) ([]models.RecordedSession, int, error) {
	sessionRecord := make([]models.RecordedSession, 0)

	query := []bson.M{
		{
			"$match": bson.M{"uid": uid},
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
	cursor, err := s.db.Collection("recorded_sessions").Aggregate(ctx, query)
	if err != nil {
		return sessionRecord, 0, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		record := new(models.RecordedSession)
		err = cursor.Decode(&record)
		if err != nil {
			return sessionRecord, 0, err
		}

		sessionRecord = append(sessionRecord, *record)
	}

	if tenant := apicontext.TenantFromContext(ctx); tenant != nil {
		query = append(query, bson.M{
			"$match": bson.M{
				"tenant_id": tenant.ID,
			},
		})
	}

	query = append(query, bson.M{
		"$count": "count",
	})

	count, err := aggregateCount(ctx, s.db.Collection("recorded_sessions"), query)
	if err != nil {
		return nil, 0, err
	}

	return sessionRecord, count, nil
}
